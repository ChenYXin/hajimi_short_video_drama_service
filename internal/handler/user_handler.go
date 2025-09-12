package handler

import (
	"net/http"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	*BaseHandler
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(),
		userService: userService,
	}
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的资料信息
// @Tags 用户
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 401 {object} models.APIResponse
// @Router /api/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := h.GetUserIDFromContext(c)
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		h.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	h.SuccessResponse(c, user)
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前登录用户的资料信息
// @Tags 用户
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UpdateProfileRequest true "更新信息"
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := h.GetUserIDFromContext(c)
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}

	var req models.UpdateProfileRequest
	if err := h.ValidateRequest(c, &req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	user, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "资料更新成功", user)
}