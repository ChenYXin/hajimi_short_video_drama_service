package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator 全局验证器实例
var Validator *validator.Validate

// ValidationError 验证错误结构
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrors 验证错误列表
type ValidationErrors []ValidationError

// Error 实现 error 接口
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// InitValidator 初始化验证器
func InitValidator() {
	Validator = validator.New()
	
	// 注册自定义标签名称
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	// 注册自定义验证规则
	registerCustomValidations()
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) ValidationErrors {
	var validationErrors ValidationErrors
	
	err := Validator.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   fmt.Sprintf("%v", err.Value()),
				Message: getErrorMessage(err),
			})
		}
	}
	
	return validationErrors
}

// getErrorMessage 获取错误消息
func getErrorMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()
	
	switch tag {
	case "required":
		return fmt.Sprintf("%s 是必填字段", field)
	case "email":
		return fmt.Sprintf("%s 必须是有效的邮箱地址", field)
	case "min":
		return fmt.Sprintf("%s 最小长度为 %s", field, param)
	case "max":
		return fmt.Sprintf("%s 最大长度为 %s", field, param)
	case "len":
		return fmt.Sprintf("%s 长度必须为 %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s 必须是以下值之一: %s", field, param)
	default:
		return fmt.Sprintf("%s 验证失败", field)
	}
}

// registerCustomValidations 注册自定义验证规则
func registerCustomValidations() {
	// 可以在这里添加自定义验证规则
	// 例如：验证手机号格式、验证视频文件格式等
}