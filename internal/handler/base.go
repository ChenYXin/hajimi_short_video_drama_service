package handler

import (
	"net/http"
	"strconv"

	"gin-mysql-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BaseHandler 基础处理器
type BaseHandler struct {
	validator *validator.Validate
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		validator: validator.New(),
	}
}

// SuccessResponse 成功响应
func (h *BaseHandler) SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "操作成功",
		Data:    data,
	})
}

// SuccessResponseWithMessage 带消息的成功响应
func (h *BaseHandler) SuccessResponseWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 错误响应
func (h *BaseHandler) ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, models.APIResponse{
		Success: false,
		Message: message,
		Error:   message,
	})
}

// ValidationErrorResponse 验证错误响应
func (h *BaseHandler) ValidationErrorResponse(c *gin.Context, err error) {
	var errorMessages []string
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errorMessages = append(errorMessages, h.getValidationErrorMessage(fieldError))
		}
	} else {
		errorMessages = append(errorMessages, err.Error())
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": "参数验证失败",
		"errors":  errorMessages,
	})
}

// ValidateRequest 验证请求参数
func (h *BaseHandler) ValidateRequest(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}
	return h.validator.Struct(req)
}

// GetPaginationParams 获取分页参数
func (h *BaseHandler) GetPaginationParams(c *gin.Context) (page, pageSize int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	return page, pageSize
}

// GetUserIDFromContext 从上下文获取用户ID
func (h *BaseHandler) GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	
	if id, ok := userID.(uint); ok {
		return id, true
	}
	
	return 0, false
}

// GetUserRoleFromContext 从上下文获取用户角色
func (h *BaseHandler) GetUserRoleFromContext(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	
	if roleStr, ok := role.(string); ok {
		return roleStr, true
	}
	
	return "", false
}

// getValidationErrorMessage 获取验证错误消息
func (h *BaseHandler) getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " 是必填字段"
	case "email":
		return fe.Field() + " 必须是有效的邮箱地址"
	case "min":
		return fe.Field() + " 长度不能少于 " + fe.Param() + " 个字符"
	case "max":
		return fe.Field() + " 长度不能超过 " + fe.Param() + " 个字符"
	case "len":
		return fe.Field() + " 长度必须是 " + fe.Param() + " 个字符"
	case "oneof":
		return fe.Field() + " 必须是以下值之一: " + fe.Param()
	default:
		return fe.Field() + " 验证失败"
	}
}