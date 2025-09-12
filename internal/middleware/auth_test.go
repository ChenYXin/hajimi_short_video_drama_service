package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gin-mysql-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	
	// 创建测试路由
	router := gin.New()
	router.Use(AuthMiddleware(jwtManager))
	router.GET("/protected", func(c *gin.Context) {
		userID := c.GetUint("user_id")
		username := c.GetString("username")
		role := c.GetString("role")
		
		c.JSON(http.StatusOK, gin.H{
			"user_id":  userID,
			"username": username,
			"role":     role,
		})
	})

	t.Run("有效token访问成功", func(t *testing.T) {
		// 生成有效 token
		token, err := jwtManager.GenerateToken(1, "testuser", "user")
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("缺少Authorization头", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("无效的Authorization格式", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("无效token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAdminAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	
	// 创建测试路由
	router := gin.New()
	router.Use(AdminAuthMiddleware(jwtManager))
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
	})

	t.Run("管理员token访问成功", func(t *testing.T) {
		// 生成管理员 token
		token, err := jwtManager.GenerateToken(1, "admin", "admin")
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("普通用户token访问被拒绝", func(t *testing.T) {
		// 生成普通用户 token
		token, err := jwtManager.GenerateToken(1, "user", "user")
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestOptionalAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	
	// 创建测试路由
	router := gin.New()
	router.Use(OptionalAuthMiddleware(jwtManager))
	router.GET("/optional", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if exists {
			c.JSON(http.StatusOK, gin.H{
				"authenticated": true,
				"user_id":      userID,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"authenticated": false,
			})
		}
	})

	t.Run("有token时设置用户信息", func(t *testing.T) {
		token, err := jwtManager.GenerateToken(1, "testuser", "user")
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/optional", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"authenticated":true`)
	})

	t.Run("无token时正常访问", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/optional", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"authenticated":false`)
	})
}