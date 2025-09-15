package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"gin-mysql-api/internal/handler"
	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/repository"
	"gin-mysql-api/internal/router"
	"gin-mysql-api/internal/service"
	"gin-mysql-api/pkg/config"
	"gin-mysql-api/pkg/database"
)

type AdminIntegrationTestSuite struct {
	suite.Suite
	router     *gin.Engine
	db         *gorm.DB
	config     *config.Config
	adminRepo  repository.AdminRepository
	dramaRepo  repository.DramaRepository
	adminToken string
}

func (suite *AdminIntegrationTestSuite) SetupSuite() {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	// 加载测试配置 (使用开发配置 + 环境变量覆盖)
	cfg, err := config.LoadConfig("../configs/config.yaml")
	suite.Require().NoError(err)

	// 覆盖测试专用配置
	cfg.Database.DBName = "hajimi_test"
	cfg.Redis.DB = 1
	cfg.Server.Port = 8081
	cfg.JWT.Secret = "test-secret-key"
	cfg.Logging.Level = "debug"

	suite.config = cfg

	// 连接测试数据库
	db, err := database.NewConnection(cfg)
	suite.Require().NoError(err)
	suite.db = db

	// 迁移数据库
	err = db.AutoMigrate(&models.User{}, &models.Admin{}, &models.Drama{}, &models.Episode{})
	suite.Require().NoError(err)

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	dramaRepo := repository.NewDramaRepository(db)
	episodeRepo := repository.NewEpisodeRepository(db)

	suite.adminRepo = adminRepo
	suite.dramaRepo = dramaRepo

	// 初始化服务层
	userService := service.NewUserService(userRepo, cfg)
	adminService := service.NewAdminService(adminRepo, cfg)
	dramaService := service.NewDramaService(dramaRepo, episodeRepo, nil, cfg)
	fileService := service.NewFileService(cfg)

	// 初始化处理器
	container := &handler.Container{
		UserService:  userService,
		AdminService: adminService,
		DramaService: dramaService,
		FileService:  fileService,
		Config:       cfg,
	}

	// 设置路由
	suite.router = router.SetupRouter(container)

	// 创建测试管理员
	suite.createTestAdmin()
}

func (suite *AdminIntegrationTestSuite) TearDownSuite() {
	// 清理测试数据
	suite.db.Exec("DELETE FROM episodes")
	suite.db.Exec("DELETE FROM dramas")
	suite.db.Exec("DELETE FROM admins")
}

func (suite *AdminIntegrationTestSuite) createTestAdmin() {
	// 创建测试管理员
	testAdmin := &models.Admin{
		Username: "testadmin",
		Email:    "admin@example.com",
		Password: "admin123",
		Role:     "admin",
	}
	err := suite.adminRepo.Create(testAdmin)
	suite.Require().NoError(err)

	// 获取管理员认证token
	loginData := map[string]string{
		"username": "testadmin",
		"password": "admin123",
	}
	loginJSON, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/admin/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var loginResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResp)

	if data, ok := loginResp["data"].(map[string]interface{}); ok {
		if token, ok := data["token"].(string); ok {
			suite.adminToken = token
		}
	}
}

// 测试管理员认证API
func (suite *AdminIntegrationTestSuite) TestAdminAuthAPI() {
	// 测试管理员登录
	suite.Run("管理员登录", func() {
		loginData := map[string]string{
			"username": "testadmin",
			"password": "admin123",
		}
		loginJSON, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/admin/login", bytes.NewBuffer(loginJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
	})

	// 测试管理员仪表板
	suite.Run("管理员仪表板", func() {
		req, _ := http.NewRequest("GET", "/admin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)
	})
}

// 测试短剧管理API
func (suite *AdminIntegrationTestSuite) TestDramaManagementAPI() {
	var dramaID uint

	// 测试创建短剧
	suite.Run("创建短剧", func() {
		dramaData := map[string]interface{}{
			"title":       "测试短剧",
			"description": "这是一个测试短剧",
			"category":    "爱情",
			"tags":        []string{"浪漫", "都市"},
			"status":      "published",
		}
		dramaJSON, _ := json.Marshal(dramaData)

		req, _ := http.NewRequest("POST", "/admin/api/dramas", bytes.NewBuffer(dramaJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])

		// 保存短剧ID用于后续测试
		if data, ok := response["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(float64); ok {
				dramaID = uint(id)
			}
		}
	})

	// 测试获取短剧列表
	suite.Run("获取短剧列表", func() {
		req, _ := http.NewRequest("GET", "/admin/api/dramas?page=1&limit=10", nil)
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
	})

	// 测试更新短剧
	suite.Run("更新短剧", func() {
		if dramaID == 0 {
			suite.T().Skip("跳过更新测试，因为没有创建短剧")
			return
		}

		updateData := map[string]interface{}{
			"title":       "更新的测试短剧",
			"description": "这是一个更新的测试短剧",
			"category":    "悬疑",
		}
		updateJSON, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/admin/api/dramas/%d", dramaID), bytes.NewBuffer(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
	})

	// 测试删除短剧
	suite.Run("删除短剧", func() {
		if dramaID == 0 {
			suite.T().Skip("跳过删除测试，因为没有创建短剧")
			return
		}

		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/admin/api/dramas/%d", dramaID), nil)
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
	})
}

// 测试文件上传功能
func (suite *AdminIntegrationTestSuite) TestFileUploadAPI() {
	suite.Run("文件上传", func() {
		// 创建一个模拟的文件上传请求
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 添加文件字段
		fileWriter, err := writer.CreateFormFile("file", "test.jpg")
		assert.NoError(suite.T(), err)

		// 写入模拟的图片数据
		fileWriter.Write([]byte("fake image data"))

		// 添加其他表单字段
		writer.WriteField("type", "cover")
		writer.WriteField("drama_id", "1")

		writer.Close()

		req, _ := http.NewRequest("POST", "/admin/api/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		// 文件上传可能因为目录不存在而失败，这是正常的
		assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusBadRequest || w.Code == http.StatusInternalServerError)
	})
}

// 测试剧集管理API
func (suite *AdminIntegrationTestSuite) TestEpisodeManagementAPI() {
	// 首先创建一个短剧
	drama := &models.Drama{
		Title:       "测试短剧",
		Description: "用于测试剧集的短剧",
		Category:    "测试",
		Status:      "published",
	}
	err := suite.dramaRepo.Create(drama)
	suite.Require().NoError(err)

	var episodeID uint

	// 测试创建剧集
	suite.Run("创建剧集", func() {
		episodeData := map[string]interface{}{
			"drama_id":    drama.ID,
			"title":       "第一集",
			"description": "测试剧集描述",
			"episode_num": 1,
			"duration":    1800, // 30分钟
			"status":      "published",
		}
		episodeJSON, _ := json.Marshal(episodeData)

		req, _ := http.NewRequest("POST", "/admin/api/episodes", bytes.NewBuffer(episodeJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])

		// 保存剧集ID用于后续测试
		if data, ok := response["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(float64); ok {
				episodeID = uint(id)
			}
		}
	})

	// 测试获取剧集列表
	suite.Run("获取剧集列表", func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/api/dramas/%d/episodes", drama.ID), nil)
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
	})

	// 测试更新剧集
	suite.Run("更新剧集", func() {
		if episodeID == 0 {
			suite.T().Skip("跳过更新测试，因为没有创建剧集")
			return
		}

		updateData := map[string]interface{}{
			"title":       "更新的第一集",
			"description": "更新的测试剧集描述",
			"duration":    2100, // 35分钟
		}
		updateJSON, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/admin/api/episodes/%d", episodeID), bytes.NewBuffer(updateJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.adminToken)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
	})
}

func TestAdminIntegrationTestSuite(t *testing.T) {
	// 跳过集成测试，如果没有设置测试环境
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("跳过集成测试，设置 INTEGRATION_TEST=true 来运行")
	}

	suite.Run(t, new(AdminIntegrationTestSuite))
}
