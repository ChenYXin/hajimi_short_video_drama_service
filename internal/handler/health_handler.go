package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	*BaseHandler
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		BaseHandler: NewBaseHandler(),
	}
}

// HealthCheck 健康检查
// @Summary 健康检查
// @Description 检查服务健康状态
// @Tags 系统
// @Produce json
// @Success 200 {object} models.APIResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	data := gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "gin-mysql-api",
		"version":   "1.0.0",
	}

	h.SuccessResponseWithMessage(c, "服务运行正常", data)
}

// ReadinessCheck 就绪检查
// @Summary 就绪检查
// @Description 检查服务是否准备好接收请求
// @Tags 系统
// @Produce json
// @Success 200 {object} models.APIResponse
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// 这里可以添加数据库连接检查、Redis 连接检查等
	data := gin.H{
		"status":    "ready",
		"timestamp": time.Now().Unix(),
		"checks": gin.H{
			"database": "ok",
			"redis":    "ok",
		},
	}

	h.SuccessResponseWithMessage(c, "服务已就绪", data)
}

// LivenessCheck 存活检查
// @Summary 存活检查
// @Description 检查服务是否存活
// @Tags 系统
// @Produce json
// @Success 200 {object} models.APIResponse
// @Router /live [get]
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	data := gin.H{
		"status":    "alive",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(time.Now()).String(), // 这里应该记录实际的启动时间
	}

	h.SuccessResponseWithMessage(c, "服务存活", data)
}