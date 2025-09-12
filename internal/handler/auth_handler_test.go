package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService 模拟认证服务
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) RegisterUser(req models.RegisterRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) LoginUser(req models.LoginRequest) (*models.LoginResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) LoginAdmin(req models.AdminLoginRequest) (*models.LoginResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(tokenString string) (string, error) {
	args := m.Called(tokenString)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) VerifyToken(tokenString string) (*utils.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.JWTClaims), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockAuthService := new(MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	t.Run("成功注册", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		user := &models.User{
			ID:       1,
			Username: req.Username,
			Email:    req.Email,
			IsActive: true,
		}

		mockAuthService.On("RegisterUser", req).Return(user, nil)

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		authHandler.Register(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "注册成功", response.Message)

		mockAuthService.AssertExpectations(t)
	})

	t.Run("注册失败 - 用户名已存在", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "existinguser",
			Email:    "test@example.com",
			Password: "password123",
		}

		mockAuthService.On("RegisterUser", req).Return((*models.User)(nil), assert.AnError)

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		authHandler.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockAuthService.AssertExpectations(t)
	})

	t.Run("注册失败 - 参数验证错误", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "", // 必填字段为空
			Email:    "invalid-email", // 无效邮箱
			Password: "123", // 密码太短
		}

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		authHandler.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockAuthService := new(MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	t.Run("成功登录", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		loginResponse := &models.LoginResponse{
			Token: "test-jwt-token",
			User: &models.User{
				ID:       1,
				Username: "testuser",
				Email:    req.Email,
			},
		}

		mockAuthService.On("LoginUser", req).Return(loginResponse, nil)

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		authHandler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "登录成功", response.Message)

		mockAuthService.AssertExpectations(t)
	})

	t.Run("登录失败 - 凭据错误", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		mockAuthService.On("LoginUser", req).Return((*models.LoginResponse)(nil), assert.AnError)

		reqBody, _ := json.Marshal(req)
		httpReq := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
		httpReq.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httpReq

		authHandler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		mockAuthService.AssertExpectations(t)
	})
}