package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-mysql-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// HTTPTestHelper HTTP 测试助手
type HTTPTestHelper struct {
	t      *testing.T
	router *gin.Engine
}

// NewHTTPTestHelper 创建 HTTP 测试助手
func NewHTTPTestHelper(t *testing.T, router *gin.Engine) *HTTPTestHelper {
	gin.SetMode(gin.TestMode)
	return &HTTPTestHelper{
		t:      t,
		router: router,
	}
}

// GET 发送 GET 请求
func (h *HTTPTestHelper) GET(path string, headers ...map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", path, nil)
	
	// 设置请求头
	for _, headerMap := range headers {
		for key, value := range headerMap {
			req.Header.Set(key, value)
		}
	}
	
	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)
	return w
}

// POST 发送 POST 请求
func (h *HTTPTestHelper) POST(path string, body interface{}, headers ...map[string]string) *httptest.ResponseRecorder {
	var reqBody io.Reader
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		assert.NoError(h.t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	}
	
	req := httptest.NewRequest("POST", path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	
	// 设置请求头
	for _, headerMap := range headers {
		for key, value := range headerMap {
			req.Header.Set(key, value)
		}
	}
	
	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)
	return w
}

// PUT 发送 PUT 请求
func (h *HTTPTestHelper) PUT(path string, body interface{}, headers ...map[string]string) *httptest.ResponseRecorder {
	var reqBody io.Reader
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		assert.NoError(h.t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	}
	
	req := httptest.NewRequest("PUT", path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	
	// 设置请求头
	for _, headerMap := range headers {
		for key, value := range headerMap {
			req.Header.Set(key, value)
		}
	}
	
	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)
	return w
}

// DELETE 发送 DELETE 请求
func (h *HTTPTestHelper) DELETE(path string, headers ...map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("DELETE", path, nil)
	
	// 设置请求头
	for _, headerMap := range headers {
		for key, value := range headerMap {
			req.Header.Set(key, value)
		}
	}
	
	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)
	return w
}

// AssertStatusCode 断言状态码
func (h *HTTPTestHelper) AssertStatusCode(w *httptest.ResponseRecorder, expectedCode int) {
	assert.Equal(h.t, expectedCode, w.Code, "Status code mismatch. Response body: %s", w.Body.String())
}

// AssertJSONResponse 断言 JSON 响应
func (h *HTTPTestHelper) AssertJSONResponse(w *httptest.ResponseRecorder, expected interface{}) {
	var actual interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	assert.NoError(h.t, err)
	assert.Equal(h.t, expected, actual)
}

// AssertAPIResponse 断言 API 响应格式
func (h *HTTPTestHelper) AssertAPIResponse(w *httptest.ResponseRecorder, success bool, message string) *models.APIResponse {
	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(h.t, err)
	assert.Equal(h.t, success, response.Success)
	if message != "" {
		assert.Equal(h.t, message, response.Message)
	}
	return &response
}

// AssertSuccessResponse 断言成功响应
func (h *HTTPTestHelper) AssertSuccessResponse(w *httptest.ResponseRecorder, message string) *models.APIResponse {
	h.AssertStatusCode(w, http.StatusOK)
	return h.AssertAPIResponse(w, true, message)
}

// AssertErrorResponse 断言错误响应
func (h *HTTPTestHelper) AssertErrorResponse(w *httptest.ResponseRecorder, statusCode int, message string) *models.APIResponse {
	h.AssertStatusCode(w, statusCode)
	return h.AssertAPIResponse(w, false, message)
}

// GetAuthHeader 获取认证头
func GetAuthHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// AssertContains 断言字符串包含
func AssertContains(t *testing.T, haystack, needle string) {
	assert.Contains(t, haystack, needle)
}

// AssertNotContains 断言字符串不包含
func AssertNotContains(t *testing.T, haystack, needle string) {
	assert.NotContains(t, haystack, needle)
}

// AssertEmpty 断言为空
func AssertEmpty(t *testing.T, obj interface{}) {
	assert.Empty(t, obj)
}

// AssertNotEmpty 断言不为空
func AssertNotEmpty(t *testing.T, obj interface{}) {
	assert.NotEmpty(t, obj)
}

// AssertEqual 断言相等
func AssertEqual(t *testing.T, expected, actual interface{}) {
	assert.Equal(t, expected, actual)
}

// AssertNotEqual 断言不相等
func AssertNotEqual(t *testing.T, expected, actual interface{}) {
	assert.NotEqual(t, expected, actual)
}

// AssertNil 断言为 nil
func AssertNil(t *testing.T, obj interface{}) {
	assert.Nil(t, obj)
}

// AssertNotNil 断言不为 nil
func AssertNotNil(t *testing.T, obj interface{}) {
	assert.NotNil(t, obj)
}

// AssertTrue 断言为 true
func AssertTrue(t *testing.T, value bool) {
	assert.True(t, value)
}

// AssertFalse 断言为 false
func AssertFalse(t *testing.T, value bool) {
	assert.False(t, value)
}

// AssertNoError 断言无错误
func AssertNoError(t *testing.T, err error) {
	assert.NoError(t, err)
}

// AssertError 断言有错误
func AssertError(t *testing.T, err error) {
	assert.Error(t, err)
}