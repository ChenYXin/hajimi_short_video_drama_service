package handler

import (
	"net/http"
	"strconv"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/service"

	"github.com/gin-gonic/gin"
)

// WebHandler Web 管理界面处理器
type WebHandler struct {
	*BaseHandler
	adminService service.AdminService
	userService  service.UserService
	dramaService service.DramaService
}

// NewWebHandler 创建 Web 处理器
func NewWebHandler(
	adminService service.AdminService,
	userService service.UserService,
	dramaService service.DramaService,
) *WebHandler {
	return &WebHandler{
		BaseHandler:  NewBaseHandler(),
		adminService: adminService,
		userService:  userService,
		dramaService: dramaService,
	}
}

// LoginPage 显示登录页面
func (h *WebHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "管理员登录",
	})
}

// Login 处理登录请求
func (h *WebHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"Title":        "管理员登录",
			"ErrorMessage": "用户名和密码不能为空",
		})
		return
	}

	// 调用管理员登录服务
	loginReq := models.AdminLoginRequest{
		Username: username,
		Password: password,
	}

	response, err := h.adminService.Login(loginReq)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Title":        "管理员登录",
			"ErrorMessage": err.Error(),
		})
		return
	}

	// 设置 session 或 cookie
	c.SetCookie("admin_token", response.Token, 86400, "/", "", false, true)

	// 重定向到仪表板
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// Logout 退出登录
func (h *WebHandler) Logout(c *gin.Context) {
	c.SetCookie("admin_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/admin/login")
}

// Dashboard 仪表板页面
func (h *WebHandler) Dashboard(c *gin.Context) {
	// 获取统计数据
	stats := h.getDashboardStats()
	
	// 获取最新短剧
	recentDramas, _ := h.dramaService.GetDramas(1, 5, "")
	
	// 获取最新用户
	recentUsers, _ := h.userService.GetUserList(1, 5)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title":        "仪表板",
		"ActiveMenu":   "dashboard",
		"AdminUser":    h.getAdminUser(c),
		"Breadcrumbs":  h.getBreadcrumbs("仪表板", ""),
		"Stats":        stats,
		"RecentDramas": recentDramas.Dramas,
		"RecentUsers":  recentUsers.Users,
	})
}

// DramasPage 短剧管理页面
func (h *WebHandler) DramasPage(c *gin.Context) {
	page, pageSize := h.GetPaginationParams(c)
	search := c.Query("search")
	genre := c.Query("genre")
	status := c.Query("status")

	dramas, err := h.adminService.GetDramaList(page, pageSize)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "dramas.html", gin.H{
			"Title":        "短剧管理",
			"ActiveMenu":   "dramas",
			"AdminUser":    h.getAdminUser(c),
			"Breadcrumbs":  h.getBreadcrumbs("短剧管理", ""),
			"ErrorMessage": "获取短剧列表失败",
		})
		return
	}

	c.HTML(http.StatusOK, "dramas.html", gin.H{
		"Title":        "短剧管理",
		"ActiveMenu":   "dramas",
		"AdminUser":    h.getAdminUser(c),
		"Breadcrumbs":  h.getBreadcrumbs("短剧管理", ""),
		"Dramas":       dramas,
		"SearchQuery":  search,
		"GenreFilter":  genre,
		"StatusFilter": status,
		"PageNumbers":  h.generatePageNumbers(dramas.Page, dramas.TotalPages),
	})
}

// EpisodesPage 剧集管理页面
func (h *WebHandler) EpisodesPage(c *gin.Context) {
	page, pageSize := h.GetPaginationParams(c)
	dramaIDStr := c.Query("drama_id")
	search := c.Query("search")

	var dramaID uint
	var drama *models.Drama
	var episodes *models.PaginatedEpisodes

	if dramaIDStr != "" {
		id, parseErr := strconv.ParseUint(dramaIDStr, 10, 32)
		if parseErr == nil {
			dramaID = uint(id)
			episodes, _ = h.adminService.GetEpisodeList(dramaID, page, pageSize)
		}
	}

	// 获取所有短剧用于选择器
	allDramas, _ := h.adminService.GetDramaList(1, 1000)

	c.HTML(http.StatusOK, "episodes.html", gin.H{
		"Title":           "剧集管理",
		"ActiveMenu":      "episodes",
		"AdminUser":       h.getAdminUser(c),
		"Breadcrumbs":     h.getBreadcrumbs("剧集管理", ""),
		"Drama":           drama,
		"Episodes":        episodes,
		"AllDramas":       allDramas.Dramas,
		"SelectedDramaID": dramaID,
		"SearchQuery":     search,
		"PageNumbers":     h.generatePageNumbers(episodes.Page, episodes.TotalPages),
	})
}

// UsersPage 用户管理页面
func (h *WebHandler) UsersPage(c *gin.Context) {
	page, pageSize := h.GetPaginationParams(c)
	search := c.Query("search")
	status := c.Query("status")
	sort := c.Query("sort")

	users, err := h.userService.GetUserList(page, pageSize)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "users.html", gin.H{
			"Title":        "用户管理",
			"ActiveMenu":   "users",
			"AdminUser":    h.getAdminUser(c),
			"Breadcrumbs":  h.getBreadcrumbs("用户管理", ""),
			"ErrorMessage": "获取用户列表失败",
		})
		return
	}

	c.HTML(http.StatusOK, "users.html", gin.H{
		"Title":        "用户管理",
		"ActiveMenu":   "users",
		"AdminUser":    h.getAdminUser(c),
		"Breadcrumbs":  h.getBreadcrumbs("用户管理", ""),
		"Users":        users,
		"SearchQuery":  search,
		"StatusFilter": status,
		"SortBy":       sort,
		"PageNumbers":  h.generatePageNumbers(users.Page, users.TotalPages),
	})
}

// 辅助方法

// getAdminUser 获取当前管理员用户信息
func (h *WebHandler) getAdminUser(c *gin.Context) gin.H {
	// 这里应该从 JWT token 或 session 中获取管理员信息
	// 暂时返回模拟数据
	return gin.H{
		"ID":       1,
		"Username": "admin",
		"Email":    "admin@example.com",
	}
}

// getBreadcrumbs 生成面包屑导航
func (h *WebHandler) getBreadcrumbs(current, parent string) []gin.H {
	breadcrumbs := []gin.H{
		{"Name": "首页", "URL": "/admin/dashboard", "Active": false},
	}

	if parent != "" {
		breadcrumbs = append(breadcrumbs, gin.H{
			"Name": parent, "URL": "", "Active": false,
		})
	}

	breadcrumbs = append(breadcrumbs, gin.H{
		"Name": current, "URL": "", "Active": true,
	})

	return breadcrumbs
}

// getDashboardStats 获取仪表板统计数据
func (h *WebHandler) getDashboardStats() gin.H {
	// 这里应该调用实际的统计服务
	// 暂时返回模拟数据
	return gin.H{
		"TotalDramas":   150,
		"TotalEpisodes": 1200,
		"TotalUsers":    5000,
		"TodayViews":    8500,
	}
}

// generatePageNumbers 生成分页页码
func (h *WebHandler) generatePageNumbers(currentPage, totalPages int) []int {
	var pages []int
	
	// 显示当前页前后各2页
	start := currentPage - 2
	if start < 1 {
		start = 1
	}
	
	end := currentPage + 2
	if end > totalPages {
		end = totalPages
	}
	
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}
	
	return pages
}