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

// MockUserRepository 模拟用户仓库
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List(offset, limit int) ([]models.User, int64, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

func TestUserService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	userService := NewUserService(mockRepo, jwtManager)

	t.Run("成功注册用户", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			Phone:    "12345678901",
		}

		// 设置 mock 期望
		mockRepo.On("ExistsByUsername", req.Username).Return(false, nil)
		mockRepo.On("ExistsByEmail", req.Email).Return(false, nil)
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		user, err := userService.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Username, user.Username)
		assert.Equal(t, req.Email, user.Email)
		assert.Empty(t, user.Password) // 密码应该被清除
		assert.True(t, user.IsActive)

		mockRepo.AssertExpectations(t)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "existinguser",
			Email:    "test@example.com",
			Password: "password123",
		}

		mockRepo.On("ExistsByUsername", req.Username).Return(true, nil)

		user, err := userService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "用户名已存在")

		mockRepo.AssertExpectations(t)
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		req := models.RegisterRequest{
			Username: "newuser",
			Email:    "existing@example.com",
			Password: "password123",
		}

		mockRepo.On("ExistsByUsername", req.Username).Return(false, nil)
		mockRepo.On("ExistsByEmail", req.Email).Return(true, nil)

		user, err := userService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "邮箱已被注册")

		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	userService := NewUserService(mockRepo, jwtManager)

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

		mockRepo.On("GetByEmail", req.Email).Return(user, nil)

		response, err := userService.Login(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		assert.NotNil(t, response.User)

		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		req := models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		mockRepo.On("GetByEmail", req.Email).Return((*models.User)(nil), errors.New("用户不存在"))

		response, err := userService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "用户名或密码错误")

		mockRepo.AssertExpectations(t)
	})

	t.Run("密码错误", func(t *testing.T) {
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

		mockRepo.On("GetByEmail", req.Email).Return(user, nil)

		response, err := userService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "用户名或密码错误")

		mockRepo.AssertExpectations(t)
	})

	t.Run("用户已被禁用", func(t *testing.T) {
		// 重新创建 mock 以避免之前的调用影响
		mockRepo := new(MockUserRepository)
		userService := NewUserService(mockRepo, jwtManager)
		
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

		mockRepo.On("GetByEmail", req.Email).Return(user, nil)

		response, err := userService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "用户账户已被禁用")

		mockRepo.AssertExpectations(t)
	})
}