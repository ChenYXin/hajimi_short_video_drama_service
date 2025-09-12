package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager(t *testing.T) {
	secretKey := "test-secret-key"
	tokenDuration := time.Hour
	manager := NewJWTManager(secretKey, tokenDuration)

	t.Run("生成和验证有效token", func(t *testing.T) {
		userID := uint(1)
		username := "testuser"
		role := "user"

		// 生成 token
		token, err := manager.GenerateToken(userID, username, role)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// 验证 token
		claims, err := manager.VerifyToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("验证无效token", func(t *testing.T) {
		invalidToken := "invalid.token.string"
		
		_, err := manager.VerifyToken(invalidToken)
		assert.Error(t, err)
	})

	t.Run("刷新token", func(t *testing.T) {
		userID := uint(2)
		username := "testuser2"
		role := "admin"

		// 生成原始 token
		originalToken, err := manager.GenerateToken(userID, username, role)
		assert.NoError(t, err)

		// 刷新 token
		newToken, err := manager.RefreshToken(originalToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, newToken)

		// 验证新 token
		claims, err := manager.VerifyToken(newToken)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("过期token验证", func(t *testing.T) {
		// 创建一个很短过期时间的管理器
		shortManager := NewJWTManager(secretKey, time.Millisecond)
		
		token, err := shortManager.GenerateToken(1, "test", "user")
		assert.NoError(t, err)

		// 等待 token 过期
		time.Sleep(time.Millisecond * 10)

		_, err = shortManager.VerifyToken(token)
		assert.Error(t, err)
	})
}