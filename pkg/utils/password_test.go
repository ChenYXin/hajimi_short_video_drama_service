package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("成功哈希密码", func(t *testing.T) {
		password := "testpassword123"
		
		hashedPassword, err := HashPassword(password)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword)
		assert.Greater(t, len(hashedPassword), len(password))
	})

	t.Run("空密码哈希", func(t *testing.T) {
		password := ""
		
		hashedPassword, err := HashPassword(password)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("相同密码产生不同哈希", func(t *testing.T) {
		password := "samepassword"
		
		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)
		
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2) // bcrypt 每次生成的哈希都不同
	})

	t.Run("长密码哈希", func(t *testing.T) {
		password := "this_is_a_long_password_with_symbols_!@#$%^&*()"
		
		hashedPassword, err := HashPassword(password)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("包含特殊字符的密码", func(t *testing.T) {
		password := "P@ssw0rd!@#$%^&*()_+-=[]{}|;:,.<>?"
		
		hashedPassword, err := HashPassword(password)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("包含中文字符的密码", func(t *testing.T) {
		password := "密码123!@#"
		
		hashedPassword, err := HashPassword(password)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})
}

func TestVerifyPassword(t *testing.T) {
	t.Run("正确密码验证成功", func(t *testing.T) {
		password := "testpassword123"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		
		isValid := VerifyPassword(hashedPassword, password)
		
		assert.True(t, isValid)
	})

	t.Run("错误密码验证失败", func(t *testing.T) {
		password := "testpassword123"
		wrongPassword := "wrongpassword"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		
		isValid := VerifyPassword(hashedPassword, wrongPassword)
		
		assert.False(t, isValid)
	})

	t.Run("空密码验证", func(t *testing.T) {
		password := ""
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		
		isValid := VerifyPassword(hashedPassword, password)
		
		assert.True(t, isValid)
	})

	t.Run("大小写敏感验证", func(t *testing.T) {
		password := "TestPassword"
		wrongCasePassword := "testpassword"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		
		isValid := VerifyPassword(hashedPassword, wrongCasePassword)
		
		assert.False(t, isValid)
	})

	t.Run("无效哈希验证失败", func(t *testing.T) {
		password := "testpassword"
		invalidHash := "invalid_hash_string"
		
		isValid := VerifyPassword(invalidHash, password)
		
		assert.False(t, isValid)
	})

	t.Run("特殊字符密码验证", func(t *testing.T) {
		password := "P@ssw0rd!@#$%^&*()"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		
		isValid := VerifyPassword(hashedPassword, password)
		
		assert.True(t, isValid)
	})

	t.Run("中文字符密码验证", func(t *testing.T) {
		password := "密码123!@#"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		
		isValid := VerifyPassword(hashedPassword, password)
		
		assert.True(t, isValid)
	})

	t.Run("长密码验证", func(t *testing.T) {
		password := "this_is_a_long_password_with_symbols_!@#$%^&*()"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		
		isValid := VerifyPassword(hashedPassword, password)
		
		assert.True(t, isValid)
	})
}

func TestPasswordHashAndVerifyIntegration(t *testing.T) {
	t.Run("完整的哈希和验证流程", func(t *testing.T) {
		testCases := []string{
			"simple",
			"Complex123!",
			"",
			"密码测试",
			"long_password_with_numbers_123_and_symbols_!@#$%^&*()",
		}

		for _, password := range testCases {
			t.Run("password: "+password, func(t *testing.T) {
				// 哈希密码
				hashedPassword, err := HashPassword(password)
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)

				// 验证正确密码
				isValid := VerifyPassword(hashedPassword, password)
				assert.True(t, isValid)

				// 验证错误密码
				wrongPassword := password + "_wrong"
				isValid = VerifyPassword(hashedPassword, wrongPassword)
				assert.False(t, isValid)
			})
		}
	})
}