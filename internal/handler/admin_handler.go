package handler

import (
	"net/http"
	"strconv"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	*BaseHandler
	adminService service.AdminService
	userService  service.UserService
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(adminService service.AdminService, userService service.UserService) *AdminHandler {
	return &AdminHandler{
		BaseHandler:  NewBaseHandler(),
		adminService: adminService,
		userService:  userService,
	}
}

// CreateDrama 创建短剧
// @Summary 创建短剧
// @Description 管理员创建新的短剧
// @Tags 管理员
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateDramaRequest true "短剧信息"
// @Success 200 {object} models.APIResponse{data=models.Drama}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /api/admin/dramas [post]
func (h *AdminHandler) CreateDrama(c *gin.Context) {
	var req models.CreateDramaRequest
	
	if err := h.ValidateRequest(c, &req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	drama, err := h.adminService.CreateDrama(req)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "短剧创建成功", drama)
}

// UpdateDrama 更新短剧
// @Summary 更新短剧
// @Description 管理员更新短剧信息
// @Tags 管理员
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "短剧ID"
// @Param request body models.UpdateDramaRequest true "更新信息"
// @Success 200 {object} models.APIResponse{data=models.Drama}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/admin/dramas/{id} [put]
func (h *AdminHandler) UpdateDrama(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的短剧ID")
		return
	}

	var req models.UpdateDramaRequest
	if err := h.ValidateRequest(c, &req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	drama, err := h.adminService.UpdateDrama(uint(id), req)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "短剧更新成功", drama)
}

// DeleteDrama 删除短剧
// @Summary 删除短剧
// @Description 管理员删除短剧
// @Tags 管理员
// @Security BearerAuth
// @Produce json
// @Param id path int true "短剧ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/admin/dramas/{id} [delete]
func (h *AdminHandler) DeleteDrama(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的短剧ID")
		return
	}

	err = h.adminService.DeleteDrama(uint(id))
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "短剧删除成功", nil)
}

// CreateEpisode 创建剧集
// @Summary 创建剧集
// @Description 管理员创建新的剧集
// @Tags 管理员
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateEpisodeRequest true "剧集信息"
// @Success 200 {object} models.APIResponse{data=models.Episode}
// @Failure 400 {object} models.APIResponse
// @Router /api/admin/episodes [post]
func (h *AdminHandler) CreateEpisode(c *gin.Context) {
	var req models.CreateEpisodeRequest
	
	if err := h.ValidateRequest(c, &req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	episode, err := h.adminService.CreateEpisode(req)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "剧集创建成功", episode)
}

// UpdateEpisode 更新剧集
// @Summary 更新剧集
// @Description 管理员更新剧集信息
// @Tags 管理员
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "剧集ID"
// @Param request body models.UpdateEpisodeRequest true "更新信息"
// @Success 200 {object} models.APIResponse{data=models.Episode}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/admin/episodes/{id} [put]
func (h *AdminHandler) UpdateEpisode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的剧集ID")
		return
	}

	var req models.UpdateEpisodeRequest
	if err := h.ValidateRequest(c, &req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	episode, err := h.adminService.UpdateEpisode(uint(id), req)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "剧集更新成功", episode)
}

// DeleteEpisode 删除剧集
// @Summary 删除剧集
// @Description 管理员删除剧集
// @Tags 管理员
// @Security BearerAuth
// @Produce json
// @Param id path int true "剧集ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/admin/episodes/{id} [delete]
func (h *AdminHandler) DeleteEpisode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的剧集ID")
		return
	}

	err = h.adminService.DeleteEpisode(uint(id))
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "剧集删除成功", nil)
}

// GetDramaList 获取短剧列表（管理员视图）
// @Summary 获取短剧列表（管理员）
// @Description 管理员获取所有短剧列表，包括未激活的
// @Tags 管理员
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.APIResponse{data=models.PaginatedDramas}
// @Router /api/admin/dramas [get]
func (h *AdminHandler) GetDramaList(c *gin.Context) {
	page, pageSize := h.GetPaginationParams(c)

	dramas, err := h.adminService.GetDramaList(page, pageSize)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "获取短剧列表失败")
		return
	}

	h.SuccessResponse(c, dramas)
}

// GetEpisodeList 获取剧集列表（管理员视图）
// @Summary 获取剧集列表（管理员）
// @Description 管理员获取指定短剧的所有剧集列表
// @Tags 管理员
// @Security BearerAuth
// @Produce json
// @Param drama_id path int true "短剧ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.APIResponse{data=models.PaginatedEpisodes}
// @Router /api/admin/dramas/{drama_id}/episodes [get]
func (h *AdminHandler) GetEpisodeList(c *gin.Context) {
	dramaIDStr := c.Param("drama_id")
	dramaID, err := strconv.ParseUint(dramaIDStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的短剧ID")
		return
	}

	page, pageSize := h.GetPaginationParams(c)

	episodes, err := h.adminService.GetEpisodeList(uint(dramaID), page, pageSize)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponse(c, episodes)
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 管理员获取用户列表
// @Tags 管理员
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.APIResponse{data=models.PaginatedUsers}
// @Router /api/admin/users [get]
func (h *AdminHandler) GetUserList(c *gin.Context) {
	page, pageSize := h.GetPaginationParams(c)

	users, err := h.userService.GetUserList(page, pageSize)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	h.SuccessResponse(c, users)
}

// ActivateUser 激活用户
// @Summary 激活用户
// @Description 管理员激活用户账户
// @Tags 管理员
// @Security BearerAuth
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/admin/users/{id}/activate [post]
func (h *AdminHandler) ActivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	err = h.userService.ActivateUser(uint(id))
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "用户激活成功", nil)
}

// DeactivateUser 禁用用户
// @Summary 禁用用户
// @Description 管理员禁用用户账户
// @Tags 管理员
// @Security BearerAuth
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/admin/users/{id}/deactivate [post]
func (h *AdminHandler) DeactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	err = h.userService.DeactivateUser(uint(id))
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "用户禁用成功", nil)
}