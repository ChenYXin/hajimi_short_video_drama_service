package service

import (
	"errors"
	"time"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/repository"
	"gin-mysql-api/pkg/utils"
)

// AuthService 认证服务接口
type AuthService interface {
	// 用户认证
	RegisterUser(req models.RegisterRequest) (*models.User, error)
	LoginUser(req models.LoginRequest) (*models.LoginResponse, error)
	
	// 管理员认证
	LoginAdmin(req models.AdminLoginRequest) (*models.LoginResponse, error)
	
	// Token 相关
	RefreshToken(tokenString string) (string, error)
	VerifyToken(tokenString string) (*utils.JWTClaims, error)
}

// authService 认证服务实现
type authService struct {
	userRepo  repository.UserRepository
	adminRepo repository.AdminRepository
	jwtManager *utils.JWTManager
}

// NewAuthService 创建新的认证服务
func NewAuthService(
	userRepo repository.UserRepository,
	adminRepo repository.AdminRepository,
	jwtManager *utils.JWTManager,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		adminRepo:  adminRepo,
		jwtManager: jwtManager,
	}
}

// RegisterUser 用户注册
func (s *authService) RegisterUser(req models.RegisterRequest) (*models.User, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, _ = s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 对密码进行哈希处理
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码处理失败")
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
		return nil, errors.New("用户创建失败")
	}

	// 清除密码字段
	user.Password = ""
	return user, nil
}

// LoginUser 用户登录
func (s *authService) LoginUser(req models.LoginRequest) (*models.LoginResponse, error) {
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
		return nil, errors.New("令牌生成失败")
	}

	// 清除密码字段
	user.Password = ""

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // 24小时后过期
		User:      user,
	}, nil
}

// LoginAdmin 管理员登录
func (s *authService) LoginAdmin(req models.AdminLoginRequest) (*models.LoginResponse, error) {
	// 根据用户名查找管理员
	admin, err := s.adminRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查管理员是否激活
	if !admin.IsActive {
		return nil, errors.New("管理员账户已被禁用")
	}

	// 验证密码
	if !utils.VerifyPassword(admin.Password, req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 JWT token
	token, err := s.jwtManager.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		return nil, errors.New("令牌生成失败")
	}

	// 清除密码字段
	admin.Password = ""

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // 24小时后过期
		User:      admin,
	}, nil
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(tokenString string) (string, error) {
	return s.jwtManager.RefreshToken(tokenString)
}

// VerifyToken 验证令牌
func (s *authService) VerifyToken(tokenString string) (*utils.JWTClaims, error) {
	return s.jwtManager.VerifyToken(tokenString)
}