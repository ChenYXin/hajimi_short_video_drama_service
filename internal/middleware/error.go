package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"gin-mysql-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "服务器内部错误",
				Error:   err,
			})
		} else if err, ok := recovered.(error); ok {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "服务器内部错误",
				Error:   err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "服务器内部错误",
				Error:   fmt.Sprintf("%v", recovered),
			})
		}
		
		// 记录错误堆栈
		fmt.Printf("Panic recovered: %v\n%s\n", recovered, debug.Stack())
		c.Abort()
	})
}

// ValidationErrorHandler 参数验证错误处理中间件
func ValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// 检查是否有验证错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			if validationErrors, ok := err.Err.(validator.ValidationErrors); ok {
				errorMessages := make([]string, 0)
				for _, fieldError := range validationErrors {
					errorMessages = append(errorMessages, getValidationErrorMessage(fieldError))
				}
				
				c.JSON(http.StatusBadRequest, models.APIResponse{
					Success: false,
					Message: "参数验证失败",
					Error:   errorMessages,
				})
				return
			}
			
			// 其他类型的错误
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "请求处理失败",
				Error:   err.Error(),
			})
		}
	}
}

// NotFoundHandler 404 错误处理
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "请求的资源不存在",
			Error:   "404 Not Found",
		})
	}
}

// MethodNotAllowedHandler 405 错误处理
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "请求方法不被允许",
			Error:   "405 Method Not Allowed",
		})
	}
}

// getValidationErrorMessage 获取验证错误消息
func getValidationErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	
	switch tag {
	case "required":
		return field + " 是必填字段"
	case "email":
		return field + " 必须是有效的邮箱地址"
	case "min":
		return field + " 长度不能少于 " + fe.Param() + " 个字符"
	case "max":
		return field + " 长度不能超过 " + fe.Param() + " 个字符"
	case "len":
		return field + " 长度必须是 " + fe.Param() + " 个字符"
	case "oneof":
		return field + " 必须是以下值之一: " + fe.Param()
	case "numeric":
		return field + " 必须是数字"
	case "alpha":
		return field + " 只能包含字母"
	case "alphanum":
		return field + " 只能包含字母和数字"
	case "url":
		return field + " 必须是有效的 URL"
	case "uuid":
		return field + " 必须是有效的 UUID"
	case "datetime":
		return field + " 必须是有效的日期时间格式"
	case "gte":
		return field + " 必须大于或等于 " + fe.Param()
	case "lte":
		return field + " 必须小于或等于 " + fe.Param()
	case "gt":
		return field + " 必须大于 " + fe.Param()
	case "lt":
		return field + " 必须小于 " + fe.Param()
	default:
		return field + " 验证失败"
	}
}