package handler

import (
	"net/http"
	"time"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	*BaseHandler
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(),
		authService: authService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册信息"
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	
	if err := h.ValidateRequest(c, &req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	user, err := h.authService.RegisterUser(req)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "注册成功", user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录信息"
// @Success 200 {object} models.APIResponse{data=models.LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	
	if err := h.ValidateRequest(c, &req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	response, err := h.authService.LoginUser(req)
	if err != nil {
		h.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// 设置过期时间
	response.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	h.SuccessResponseWithMessage(c, "登录成功", response)
}

// AdminLogin 管理员登录
// @Summary 管理员登录
// @Description 管理员登录获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.AdminLoginRequest true "管理员登录信息"
// @Success 200 {object} models.APIResponse{data=models.LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/auth/admin/login [post]
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req models.AdminLoginRequest
	
	if err := h.ValidateRequest(c, &req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	response, err := h.authService.LoginAdmin(req)
	if err != nil {
		h.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// 设置过期时间
	response.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	h.SuccessResponseWithMessage(c, "登录成功", response)
}

// RefreshToken 刷新令牌
// @Summary 刷新访问令牌
// @Description 使用当前令牌获取新的访问令牌
// @Tags 认证
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.APIResponse{data=string}
// @Failure 401 {object} models.APIResponse
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 从 Authorization header 获取当前 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.ErrorResponse(c, http.StatusUnauthorized, "缺少认证令牌")
		return
	}

	// 提取 token
	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		h.ErrorResponse(c, http.StatusUnauthorized, "无效的认证令牌格式")
		return
	}

	// 刷新 token
	newToken, err := h.authService.RefreshToken(tokenString)
	if err != nil {
		h.ErrorResponse(c, http.StatusUnauthorized, "令牌刷新失败")
		return
	}

	response := models.LoginResponse{
		Token:     newToken,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	h.SuccessResponseWithMessage(c, "令牌刷新成功", response)
}