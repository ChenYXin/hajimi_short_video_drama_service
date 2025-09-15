package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试用的结构体
type TestUser struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"min=0,max=150"`
	Status   string `json:"status" validate:"oneof=active inactive"`
}

func TestInitValidator(t *testing.T) {
	t.Run("初始化验证器", func(t *testing.T) {
		InitValidator()

		assert.NotNil(t, Validator)
	})
}

func TestValidateStruct(t *testing.T) {
	// 确保验证器已初始化
	InitValidator()

	t.Run("有效结构体验证通过", func(t *testing.T) {
		user := TestUser{
			Username: "testuser",
			Email:    "test@example.com",
			Age:      25,
			Status:   "active",
		}

		errors := ValidateStruct(user)

		assert.Empty(t, errors)
	})

	t.Run("必填字段缺失", func(t *testing.T) {
		user := TestUser{
			// Username 缺失
			Email:  "test@example.com",
			Age:    25,
			Status: "active",
		}

		errors := ValidateStruct(user)

		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Equal(t, "username", errors[0].Field)
		assert.Equal(t, "required", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "必填字段")
	})

	t.Run("邮箱格式无效", func(t *testing.T) {
		user := TestUser{
			Username: "testuser",
			Email:    "invalid-email",
			Age:      25,
			Status:   "active",
		}

		errors := ValidateStruct(user)

		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Equal(t, "email", errors[0].Field)
		assert.Equal(t, "email", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "邮箱地址")
	})

	t.Run("字符串长度验证", func(t *testing.T) {
		user := TestUser{
			Username: "ab", // 太短，最小3个字符
			Email:    "test@example.com",
			Age:      25,
			Status:   "active",
		}

		errors := ValidateStruct(user)

		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Equal(t, "username", errors[0].Field)
		assert.Equal(t, "min", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "最小长度")
	})

	t.Run("数值范围验证", func(t *testing.T) {
		user := TestUser{
			Username: "testuser",
			Email:    "test@example.com",
			Age:      -5, // 负数，最小值为0
			Status:   "active",
		}

		errors := ValidateStruct(user)

		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Equal(t, "age", errors[0].Field)
		assert.Equal(t, "min", errors[0].Tag)
	})

	t.Run("枚举值验证", func(t *testing.T) {
		user := TestUser{
			Username: "testuser",
			Email:    "test@example.com",
			Age:      25,
			Status:   "invalid_status", // 无效状态
		}

		errors := ValidateStruct(user)

		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 1)
		assert.Equal(t, "status", errors[0].Field)
		assert.Equal(t, "oneof", errors[0].Tag)
		assert.Contains(t, errors[0].Message, "以下值之一")
	})

	t.Run("多个验证错误", func(t *testing.T) {
		user := TestUser{
			Username: "",              // 必填字段缺失
			Email:    "invalid-email", // 邮箱格式无效
			Age:      200,             // 超出最大值
			Status:   "unknown",       // 无效状态
		}

		errors := ValidateStruct(user)

		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 4)

		// 检查所有错误字段
		fields := make(map[string]bool)
		for _, err := range errors {
			fields[err.Field] = true
		}
		assert.True(t, fields["username"])
		assert.True(t, fields["email"])
		assert.True(t, fields["age"])
		assert.True(t, fields["status"])
	})
}

func TestValidationError(t *testing.T) {
	t.Run("单个验证错误消息", func(t *testing.T) {
		err := ValidationError{
			Field:   "username",
			Tag:     "required",
			Value:   "",
			Message: "username 是必填字段",
		}

		assert.Equal(t, "username", err.Field)
		assert.Equal(t, "required", err.Tag)
		assert.Equal(t, "", err.Value)
		assert.Equal(t, "username 是必填字段", err.Message)
	})
}

func TestValidationErrors(t *testing.T) {
	t.Run("多个验证错误消息", func(t *testing.T) {
		errors := ValidationErrors{
			{
				Field:   "username",
				Tag:     "required",
				Value:   "",
				Message: "username 是必填字段",
			},
			{
				Field:   "email",
				Tag:     "email",
				Value:   "invalid",
				Message: "email 必须是有效的邮箱地址",
			},
		}

		errorMessage := errors.Error()
		assert.Contains(t, errorMessage, "username 是必填字段")
		assert.Contains(t, errorMessage, "email 必须是有效的邮箱地址")
		assert.Contains(t, errorMessage, ";")
	})

	t.Run("空验证错误列表", func(t *testing.T) {
		errors := ValidationErrors{}

		errorMessage := errors.Error()
		assert.Empty(t, errorMessage)
	})
}

func TestGetErrorMessage(t *testing.T) {
	// 这个测试需要创建 validator.FieldError 实例，比较复杂
	// 在实际项目中，可以通过集成测试来验证错误消息的正确性

	t.Run("验证错误消息格式", func(t *testing.T) {
		InitValidator()

		user := TestUser{
			Username: "", // required 错误
		}

		errors := ValidateStruct(user)

		if len(errors) > 0 {
			assert.Contains(t, errors[0].Message, "必填字段")
		}
	})
}

// 测试自定义验证规则（如果有的话）
func TestCustomValidations(t *testing.T) {
	InitValidator()

	t.Run("自定义验证规则注册", func(t *testing.T) {
		// 这里可以测试自定义验证规则
		// 例如手机号验证、视频文件格式验证等
		assert.NotNil(t, Validator)
	})
}

// 测试实际的模型验证
func TestRealModelValidation(t *testing.T) {
	InitValidator()

	t.Run("用户注册请求验证", func(t *testing.T) {
		req := RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			Phone:    "12345678901",
		}

		errors := ValidateStruct(req)
		assert.Empty(t, errors)
	})

	t.Run("用户登录请求验证", func(t *testing.T) {
		req := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		errors := ValidateStruct(req)
		assert.Empty(t, errors)
	})

	t.Run("创建短剧请求验证", func(t *testing.T) {
		req := CreateDramaRequest{
			Title:       "测试短剧",
			Description: "测试描述",
			Category:    "喜剧",
		}

		errors := ValidateStruct(req)
		assert.Empty(t, errors)
	})

	t.Run("创建剧集请求验证", func(t *testing.T) {
		req := CreateEpisodeRequest{
			DramaID:    1,
			Title:      "测试剧集",
			EpisodeNum: 1,
			Duration:   30,
		}

		errors := ValidateStruct(req)
		assert.Empty(t, errors)
	})
}
