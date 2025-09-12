package main

import (
	"fmt"
	"log"
	"time"

	"gin-mysql-api/internal/middleware"
	"gin-mysql-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

// 这是一个展示如何使用 JWT 认证系统的示例
func main() {
	// 1. 创建 JWT 管理器
	jwtManager := utils.NewJWTManager("your-secret-key", 24*time.Hour)

	// 2. 演示密码哈希和验证
	password := "mypassword123"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal("密码哈希失败:", err)
	}
	fmt.Printf("原始密码: %s\n", password)
	fmt.Printf("哈希后密码: %s\n", hashedPassword)

	// 验证密码
	isValid := utils.VerifyPassword(hashedPassword, password)
	fmt.Printf("密码验证结果: %t\n", isValid)

	// 3. 演示 JWT token 生成和验证
	userID := uint(1)
	username := "testuser"
	role := "user"

	// 生成 token
	token, err := jwtManager.GenerateToken(userID, username, role)
	if err != nil {
		log.Fatal("Token 生成失败:", err)
	}
	fmt.Printf("生成的 JWT Token: %s\n", token)

	// 验证 token
	claims, err := jwtManager.VerifyToken(token)
	if err != nil {
		log.Fatal("Token 验证失败:", err)
	}
	fmt.Printf("Token 验证成功 - 用户ID: %d, 用户名: %s, 角色: %s\n", 
		claims.UserID, claims.Username, claims.Role)

	// 4. 演示在 Gin 路由中使用认证中间件
	setupGinRoutes(jwtManager)
}

func setupGinRoutes(jwtManager *utils.JWTManager) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 公开路由
	r.POST("/login", func(c *gin.Context) {
		// 这里应该验证用户凭据
		// 为了演示，我们直接生成 token
		token, _ := jwtManager.GenerateToken(1, "testuser", "user")
		c.JSON(200, gin.H{
			"token": token,
		})
	})

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.GetUint("user_id")
			username := c.GetString("username")
			role := c.GetString("role")

			c.JSON(200, gin.H{
				"user_id":  userID,
				"username": username,
				"role":     role,
			})
		})
	}

	// 管理员路由
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware(jwtManager))
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "欢迎访问管理员面板",
			})
		})
	}

	// 可选认证路由
	public := r.Group("/public")
	public.Use(middleware.OptionalAuthMiddleware(jwtManager))
	{
		public.GET("/content", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if exists {
				c.JSON(200, gin.H{
					"message": "已认证用户的内容",
					"user_id": userID,
				})
			} else {
				c.JSON(200, gin.H{
					"message": "公开内容",
				})
			}
		})
	}

	fmt.Println("JWT 认证系统演示完成！")
	fmt.Println("可以使用以下路由进行测试:")
	fmt.Println("POST /login - 获取 token")
	fmt.Println("GET /api/profile - 需要用户认证")
	fmt.Println("GET /admin/dashboard - 需要管理员认证")
	fmt.Println("GET /public/content - 可选认证")
}