package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type IntegrationTestSuite struct {
	suite.Suite
	router   *gin.Engine
	db       *gorm.DB
	config   *config.Config
	userRepo repository.UserRepository
	authToken string
	adminToken string
}

func (suite *IntegrationTestSuite) SetupSuite() {
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
	cfg.Prometheus.Enabled = false
	
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
	
	suite.userRepo = userRepo
	
	// 初始化服务层
	userService := service.NewUserService(userRepo, cfg)
	adminService := service.NewAdminService(adminRepo, cfg)
	dramaService := service.NewDramaService(dramaRepo, episodeRepo, nil, cfg) // Redis为nil用于测试
	fileService := service.NewFileService(cfg)
	
	// 初始化处理器
	container := &handler.Container{
		UserService:    userService,
		AdminService:   adminService,
		DramaService:   dramaService,
		FileService:    fileService,
		Config:         cfg,
	}
	
	// 设置路由
	suite.router = router.SetupRouter(container)
	
	// 创建测试用户和管理员
	suite.createTestData()
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	// 清理测试数据
	suite.db.Exec("DELETE FROM episodes")
	suite.db.Exec("DELETE FROM dramas")
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM admins")
}

func (suite *IntegrationTestSuite) createTestData() {
	// 创建测试用户
	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	err := suite.userRepo.Create(testUser)
	suite.Require().NoError(err)
	
	// 获取用户认证token
	loginData := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	loginJSON, _ := json.Marshal(loginData)
	
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	var loginResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	
	if data, ok := loginResp["data"].(map[string]interface{}); ok {
		if token, ok := data["token"].(string); ok {
			suite.authToken = token
		}
	}
}

// 测试用户认证API
func (suite *IntegrationTestSuite) TestUserAuthAPI() {
	// 测试用户注册
	suite.Run("用户注册", func() {
		registerData := map[string]string{
			"username": "newuser",
			"email":    "newuser@example.com",
			"password": "password123",
		}
		registerJSON, _ := json.Marshal(registerData)
		
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(registerJSON))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		
		assert.Equal(suite.T(), http.StatusCreated, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
	})
	
	// 测试用户登录
	suite.Run("用户登录", func() {
		loginData := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		loginJSON, _ := json.Marshal(loginData)
		
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		
		assert.Equal(suite.T(), http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
		
		data := response["data"].(map[string]interface{})
		assert.NotEmpty(suite.T(), data["token"])
	})
	
	// 测试获取用户信息
	suite.Run("获取用户信息", func() {
		req, _ := http.NewRequest("GET", "/api/user/profile", nil)
		req.Header.Set("Authorization", "Bearer "+suite.authToken)
		
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		
		assert.Equal(suite.T(), http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
	})
}

// 测试短剧和剧集API
func (suite *IntegrationTestSuite) TestDramaAPI() {
	// 测试获取短剧列表
	suite.Run("获取短剧列表", func() {
		req, _ := http.NewRequest("GET", "/api/dramas?page=1&limit=10", nil)
		
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		
		assert.Equal(suite.T(), http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
	})
	
	// 测试搜索短剧
	suite.Run("搜索短剧", func() {
		req, _ := http.NewRequest("GET", "/api/dramas/search?keyword=test&page=1&limit=10", nil)
		
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		
		assert.Equal(suite.T(), http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "success", response["status"])
	})
	
	// 测试获取短剧详情
	suite.Run("获取短剧详情", func() {
		req, _ := http.NewRequest("GET", "/api/dramas/1", nil)
		
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		
		// 如果短剧不存在，应该返回404
		assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusNotFound)
	})
	
	// 测试获取剧集列表
	suite.Run("获取剧集列表", func() {
		req, _ := http.NewRequest("GET", "/api/dramas/1/episodes", nil)
		
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		
		// 如果短剧不存在，应该返回404
		assert.True(suite.T(), w.Code == http.StatusOK || w.Code == http.StatusNotFound)
	})
}

// 测试健康检查API
func (suite *IntegrationTestSuite) TestHealthAPI() {
	suite.Run("健康检查", func() {
		req, _ := http.NewRequest("GET", "/health", nil)
		
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		
		assert.Equal(suite.T(), http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "healthy", response["status"])
	})
}

func TestIntegrationTestSuite(t *testing.T) {
	// 跳过集成测试，如果没有设置测试环境
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("跳过集成测试，设置 INTEGRATION_TEST=true 来运行")
	}
	
	suite.Run(t, new(IntegrationTestSuite))
}