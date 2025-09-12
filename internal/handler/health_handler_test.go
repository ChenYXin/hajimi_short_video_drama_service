package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-mysql-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	healthHandler := NewHealthHandler()

	t.Run("健康检查", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		req := httptest.NewRequest("GET", "/health", nil)
		c.Request = req

		healthHandler.HealthCheck(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "服务运行正常", response.Message)
		
		// 检查返回的数据
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "ok", data["status"])
		assert.Equal(t, "gin-mysql-api", data["service"])
		assert.Equal(t, "1.0.0", data["version"])
		assert.NotNil(t, data["timestamp"])
	})

	t.Run("就绪检查", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		req := httptest.NewRequest("GET", "/ready", nil)
		c.Request = req

		healthHandler.ReadinessCheck(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "服务已就绪", response.Message)
		
		// 检查返回的数据
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "ready", data["status"])
		
		checks, ok := data["checks"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "ok", checks["database"])
		assert.Equal(t, "ok", checks["redis"])
	})

	t.Run("存活检查", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		req := httptest.NewRequest("GET", "/live", nil)
		c.Request = req

		healthHandler.LivenessCheck(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "服务存活", response.Message)
		
		// 检查返回的数据
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "alive", data["status"])
		assert.NotNil(t, data["timestamp"])
		assert.NotNil(t, data["uptime"])
	})
}