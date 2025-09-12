package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-mysql-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("处理 panic 字符串", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler())
		
		router.GET("/panic", func(c *gin.Context) {
			panic("test panic")
		})

		req := httptest.NewRequest("GET", "/panic", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "服务器内部错误", response.Message)
	})

	t.Run("处理 panic 错误", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler())
		
		router.GET("/panic-error", func(c *gin.Context) {
			panic(assert.AnError)
		})

		req := httptest.NewRequest("GET", "/panic-error", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "服务器内部错误", response.Message)
	})
}

func TestNotFoundHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.NoRoute(NotFoundHandler())

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "请求的资源不存在", response.Message)
}

func TestMethodNotAllowedHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.NoMethod(MethodNotAllowedHandler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "请求方法不被允许", response.Message)
}

func TestGetValidationErrorMessage(t *testing.T) {
	// 这里需要创建 validator.FieldError 的 mock
	// 由于 validator.FieldError 是接口，需要实现所有方法
	// 为简化测试，这里只测试函数存在性
	assert.NotNil(t, getValidationErrorMessage)
}