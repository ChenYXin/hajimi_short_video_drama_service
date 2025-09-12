package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-mysql-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBaseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	baseHandler := NewBaseHandler()

	t.Run("成功响应", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := gin.H{"message": "test"}
		baseHandler.SuccessResponse(c, data)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "操作成功", response.Message)
	})

	t.Run("错误响应", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		baseHandler.ErrorResponse(c, http.StatusBadRequest, "测试错误")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "测试错误", response.Message)
	})

	t.Run("获取分页参数", func(t *testing.T) {
		// 测试默认值
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		page, pageSize := baseHandler.GetPaginationParams(c)
		assert.Equal(t, 1, page)
		assert.Equal(t, 20, pageSize)

		// 测试自定义值
		req = httptest.NewRequest("GET", "/test?page=2&page_size=10", nil)
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = req

		page, pageSize = baseHandler.GetPaginationParams(c)
		assert.Equal(t, 2, page)
		assert.Equal(t, 10, pageSize)

		// 测试边界值
		req = httptest.NewRequest("GET", "/test?page=0&page_size=200", nil)
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = req

		page, pageSize = baseHandler.GetPaginationParams(c)
		assert.Equal(t, 1, page)      // 最小值为 1
		assert.Equal(t, 20, pageSize) // 最大值为 100，超出时使用默认值
	})

	t.Run("验证请求参数", func(t *testing.T) {
		type TestRequest struct {
			Name  string `json:"name" validate:"required"`
			Email string `json:"email" validate:"required,email"`
		}

		// 有效请求
		validReq := TestRequest{
			Name:  "test",
			Email: "test@example.com",
		}
		reqBody, _ := json.Marshal(validReq)
		
		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		var testReq TestRequest
		err := baseHandler.ValidateRequest(c, &testReq)
		assert.NoError(t, err)
		assert.Equal(t, "test", testReq.Name)
		assert.Equal(t, "test@example.com", testReq.Email)

		// 无效请求
		invalidReq := TestRequest{
			Name:  "", // 必填字段为空
			Email: "invalid-email", // 无效邮箱
		}
		reqBody, _ = json.Marshal(invalidReq)
		
		req = httptest.NewRequest("POST", "/test", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		c, _ = gin.CreateTestContext(w)
		c.Request = req

		var testReq2 TestRequest
		err = baseHandler.ValidateRequest(c, &testReq2)
		assert.Error(t, err)
	})
}