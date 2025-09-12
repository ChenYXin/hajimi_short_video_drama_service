package service

import (
	"errors"
	"fmt"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/repository"
	"gin-mysql-api/pkg/utils"
)

// UserService 用户服务接口
type UserService interface {
	Register(req models.RegisterRequest) (*models.User, error)
	Login(req models.LoginRequest) (*models.LoginResponse, error)
	GetProfile(userID uint) (*models.User, error)
	UpdateProfile(userID uint, req models.UpdateProfileRequest) (*models.User, error)
	GetUserList(page, pageSize int) (*models.PaginatedUsers, error)
	DeleteUser(userID uint) error
	ActivateUser(userID uint) error
	DeactivateUser(userID uint) error
}

// userService 用户服务实现
type userService struct {
	userRepo   repository.UserRepository
	jwtManager *utils.JWTManager
}

// NewUserService 创建新的用户服务
func NewUserService(userRepo repository.UserRepository, jwtManager *utils.JWTManager) UserService {
	return &userService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register 用户注册
func (s *userService) Register(req models.RegisterRequest) (*models.User, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return nil, errors.New("邮箱已被注册")
	}

	// 对密码进行哈希处理
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码处理失败: %w", err)
	}

	// 创建新用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
		IsActive: true,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("用户创建失败: %w", err)
	}

	// 清除密码字段
	user.Password = ""
	return user, nil
}

// Login 用户登录
func (s *userService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	// 根据邮箱查找用户
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户是否激活
	if !user.IsActive {
		return nil, errors.New("用户账户已被禁用")
	}

	// 验证密码
	if !utils.VerifyPassword(user.Password, req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		return nil, fmt.Errorf("令牌生成失败: %w", err)
	}

	// 清除密码字段
	user.Password = ""

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: 0, // 将在控制器中设置
		User:      user,
	}, nil
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 清除密码字段
	user.Password = ""
	return user, nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(userID uint, req models.UpdateProfileRequest) (*models.User, error) {
	// 获取现有用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 如果要更新用户名，检查是否已存在
	if req.Username != "" && req.Username != user.Username {
		exists, err := s.userRepo.ExistsByUsername(req.Username)
		if err != nil {
			return nil, fmt.Errorf("检查用户名失败: %w", err)
		}
		if exists {
			return nil, errors.New("用户名已存在")
		}
		user.Username = req.Username
	}

	// 更新其他字段
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	// 保存更新
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("更新用户信息失败: %w", err)
	}

	// 清除密码字段
	user.Password = ""
	return user, nil
}

// GetUserList 获取用户列表（管理员功能）
func (s *userService) GetUserList(page, pageSize int) (*models.PaginatedUsers, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	users, total, err := s.userRepo.List(offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 清除所有用户的密码字段
	for i := range users {
		users[i].Password = ""
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	return &models.PaginatedUsers{
		Users:       users,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}, nil
}

// DeleteUser 删除用户（管理员功能）
func (s *userService) DeleteUser(userID uint) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	err = s.userRepo.Delete(userID)
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}

// ActivateUser 激活用户（管理员功能）
func (s *userService) ActivateUser(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	user.IsActive = true
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("激活用户失败: %w", err)
	}

	return nil
}

// DeactivateUser 禁用用户（管理员功能）
func (s *userService) DeactivateUser(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	user.IsActive = false
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("禁用用户失败: %w", err)
	}

	return nil
}