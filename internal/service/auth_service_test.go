package service

import (
	"errors"
	"testing"
	"time"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_RegisterUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAdminRepo := new(MockAdminRepository)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)

	t.Run("成功注册用户", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			Phone:    "12345678901",
		}

		// 设置 mock 期望
		mockUserRepo.On("GetByUsername", req.Username).Return((*models.User)(nil), errors.New("用户不存在"))
		mockUserRepo.On("GetByEmail", req.Email).Return((*models.User)(nil), errors.New("用户不存在"))
		mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		user, err := authService.RegisterUser(req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Username, user.Username)
		assert.Equal(t, req.Email, user.Email)
		assert.Empty(t, user.Password) // 密码应该被清除
		assert.True(t, user.IsActive)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)
		
		req := models.RegisterRequest{
			Username: "existinguser",
			Email:    "test@example.com",
			Password: "password123",
		}

		existingUser := &models.User{Username: req.Username}
		mockUserRepo.On("GetByUsername", req.Username).Return(existingUser, nil)

		user, err := authService.RegisterUser(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "用户名已存在")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)
		
		req := models.RegisterRequest{
			Username: "newuser",
			Email:    "existing@example.com",
			Password: "password123",
		}

		existingUser := &models.User{Email: req.Email}
		mockUserRepo.On("GetByUsername", req.Username).Return((*models.User)(nil), errors.New("用户不存在"))
		mockUserRepo.On("GetByEmail", req.Email).Return(existingUser, nil)

		user, err := authService.RegisterUser(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "邮箱已被注册")

		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_LoginUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAdminRepo := new(MockAdminRepository)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)

	t.Run("成功登录", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		// 创建哈希密码
		hashedPassword, _ := utils.HashPassword(req.Password)
		user := &models.User{
			ID:       1,
			Username: "testuser",
			Email:    req.Email,
			Password: hashedPassword,
			IsActive: true,
		}

		mockUserRepo.On("GetByEmail", req.Email).Return(user, nil)

		response, err := authService.LoginUser(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		assert.NotNil(t, response.User)
		assert.Greater(t, response.ExpiresAt, int64(0))

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)
		
		req := models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		mockUserRepo.On("GetByEmail", req.Email).Return((*models.User)(nil), errors.New("用户不存在"))

		response, err := authService.LoginUser(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "用户名或密码错误")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("密码错误", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)
		
		req := models.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		hashedPassword, _ := utils.HashPassword("correctpassword")
		user := &models.User{
			ID:       1,
			Username: "testuser",
			Email:    req.Email,
			Password: hashedPassword,
			IsActive: true,
		}

		mockUserRepo.On("GetByEmail", req.Email).Return(user, nil)

		response, err := authService.LoginUser(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "用户名或密码错误")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("用户已被禁用", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)
		
		req := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		hashedPassword, _ := utils.HashPassword(req.Password)
		user := &models.User{
			ID:       1,
			Username: "testuser",
			Email:    req.Email,
			Password: hashedPassword,
			IsActive: false, // 用户被禁用
		}

		mockUserRepo.On("GetByEmail", req.Email).Return(user, nil)

		response, err := authService.LoginUser(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "用户账户已被禁用")

		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_LoginAdmin(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAdminRepo := new(MockAdminRepository)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)

	t.Run("成功登录", func(t *testing.T) {
		req := models.AdminLoginRequest{
			Username: "admin",
			Password: "password123",
		}

		hashedPassword, _ := utils.HashPassword(req.Password)
		admin := &models.Admin{
			ID:       1,
			Username: req.Username,
			Password: hashedPassword,
			IsActive: true,
		}

		mockAdminRepo.On("GetByUsername", req.Username).Return(admin, nil)

		response, err := authService.LoginAdmin(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		assert.NotNil(t, response.User)

		mockAdminRepo.AssertExpectations(t)
	})

	t.Run("管理员不存在", func(t *testing.T) {
		mockAdminRepo := new(MockAdminRepository)
		authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)
		
		req := models.AdminLoginRequest{
			Username: "nonexistent",
			Password: "password123",
		}

		mockAdminRepo.On("GetByUsername", req.Username).Return((*models.Admin)(nil), errors.New("管理员不存在"))

		response, err := authService.LoginAdmin(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "用户名或密码错误")

		mockAdminRepo.AssertExpectations(t)
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAdminRepo := new(MockAdminRepository)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)

	t.Run("成功刷新token", func(t *testing.T) {
		// 生成原始 token
		originalToken, err := jwtManager.GenerateToken(1, "testuser", "user")
		assert.NoError(t, err)

		// 刷新 token
		newToken, err := authService.RefreshToken(originalToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, newToken)
		assert.NotEqual(t, originalToken, newToken)

		// 验证新 token
		claims, err := authService.VerifyToken(newToken)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), claims.UserID)
		assert.Equal(t, "testuser", claims.Username)
		assert.Equal(t, "user", claims.Role)
	})
}

func TestAuthService_VerifyToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAdminRepo := new(MockAdminRepository)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	authService := NewAuthService(mockUserRepo, mockAdminRepo, jwtManager)

	t.Run("成功验证token", func(t *testing.T) {
		userID := uint(1)
		username := "testuser"
		role := "user"

		// 生成 token
		token, err := jwtManager.GenerateToken(userID, username, role)
		assert.NoError(t, err)

		// 验证 token
		claims, err := authService.VerifyToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("无效token", func(t *testing.T) {
		invalidToken := "invalid.token.string"

		claims, err := authService.VerifyToken(invalidToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}