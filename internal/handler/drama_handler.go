package handler

import (
	"net/http"
	"strconv"

	"gin-mysql-api/internal/service"

	"github.com/gin-gonic/gin"
)

// DramaHandler 短剧处理器
type DramaHandler struct {
	*BaseHandler
	dramaService service.DramaService
}

// NewDramaHandler 创建短剧处理器
func NewDramaHandler(dramaService service.DramaService) *DramaHandler {
	return &DramaHandler{
		BaseHandler:  NewBaseHandler(),
		dramaService: dramaService,
	}
}

// GetDramas 获取短剧列表
// @Summary 获取短剧列表
// @Description 获取短剧列表，支持分页和类型筛选
// @Tags 短剧
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category query string false "类型筛选"
// @Success 200 {object} models.APIResponse{data=models.PaginatedDramas}
// @Router /api/dramas [get]
func (h *DramaHandler) GetDramas(c *gin.Context) {
	page, pageSize := h.GetPaginationParams(c)
	category := c.Query("category")

	dramas, err := h.dramaService.GetDramas(page, pageSize, category)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "获取短剧列表失败")
		return
	}

	h.SuccessResponse(c, dramas)
}

// GetDramaByID 获取短剧详情
// @Summary 获取短剧详情
// @Description 根据ID获取短剧详细信息
// @Tags 短剧
// @Produce json
// @Param id path int true "短剧ID"
// @Success 200 {object} models.APIResponse{data=models.Drama}
// @Failure 404 {object} models.APIResponse
// @Router /api/dramas/{id} [get]
func (h *DramaHandler) GetDramaByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的短剧ID")
		return
	}

	drama, err := h.dramaService.GetDramaByID(uint(id))
	if err != nil {
		h.ErrorResponse(c, http.StatusNotFound, "短剧不存在")
		return
	}

	// 增加观看次数
	go h.dramaService.IncrementDramaViewCount(uint(id))

	h.SuccessResponse(c, drama)
}

// GetDramaWithEpisodes 获取短剧及其剧集
// @Summary 获取短剧及其剧集
// @Description 获取短剧详情以及所有剧集信息
// @Tags 短剧
// @Produce json
// @Param id path int true "短剧ID"
// @Success 200 {object} models.APIResponse{data=models.Drama}
// @Failure 404 {object} models.APIResponse
// @Router /api/dramas/{id}/episodes [get]
func (h *DramaHandler) GetDramaWithEpisodes(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的短剧ID")
		return
	}

	drama, err := h.dramaService.GetDramaWithEpisodes(uint(id))
	if err != nil {
		h.ErrorResponse(c, http.StatusNotFound, "短剧不存在")
		return
	}

	h.SuccessResponse(c, drama)
}

// GetEpisodesByDramaID 获取短剧的剧集列表
// @Summary 获取短剧的剧集列表
// @Description 分页获取指定短剧的剧集列表
// @Tags 短剧
// @Produce json
// @Param id path int true "短剧ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.APIResponse{data=models.PaginatedEpisodes}
// @Failure 404 {object} models.APIResponse
// @Router /api/dramas/{id}/episodes/list [get]
func (h *DramaHandler) GetEpisodesByDramaID(c *gin.Context) {
	idStr := c.Param("id")
	dramaID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的短剧ID")
		return
	}

	page, pageSize := h.GetPaginationParams(c)

	episodes, err := h.dramaService.GetEpisodesByDramaID(uint(dramaID), page, pageSize)
	if err != nil {
		h.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	h.SuccessResponse(c, episodes)
}

// GetEpisodeByID 获取剧集详情
// @Summary 获取剧集详情
// @Description 根据ID获取剧集详细信息和播放地址
// @Tags 剧集
// @Produce json
// @Param id path int true "剧集ID"
// @Success 200 {object} models.APIResponse{data=models.Episode}
// @Failure 404 {object} models.APIResponse
// @Router /api/episodes/{id} [get]
func (h *DramaHandler) GetEpisodeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "无效的剧集ID")
		return
	}

	episode, err := h.dramaService.GetEpisodeByID(uint(id))
	if err != nil {
		h.ErrorResponse(c, http.StatusNotFound, "剧集不存在")
		return
	}

	// 增加观看次数
	go h.dramaService.IncrementEpisodeViewCount(uint(id))

	h.SuccessResponse(c, episode)
}

// SearchDramas 搜索短剧
// @Summary 搜索短剧
// @Description 根据关键词搜索短剧
// @Tags 短剧
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.APIResponse{data=models.PaginatedDramas}
// @Router /api/dramas/search [get]
func (h *DramaHandler) SearchDramas(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		h.ErrorResponse(c, http.StatusBadRequest, "搜索关键词不能为空")
		return
	}

	page, pageSize := h.GetPaginationParams(c)

	dramas, err := h.dramaService.SearchDramas(keyword, page, pageSize)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "搜索失败")
		return
	}

	h.SuccessResponse(c, dramas)
}

// GetPopularDramas 获取热门短剧
// @Summary 获取热门短剧
// @Description 获取热门短剧列表
// @Tags 短剧
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.APIResponse{data=models.PaginatedDramas}
// @Router /api/dramas/popular [get]
func (h *DramaHandler) GetPopularDramas(c *gin.Context) {
	page, pageSize := h.GetPaginationParams(c)

	dramas, err := h.dramaService.GetPopularDramas(page, pageSize)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "获取热门短剧失败")
		return
	}

	h.SuccessResponse(c, dramas)
}
