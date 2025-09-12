package testutil

import (
	"fmt"
	"log"

	"gin-mysql-api/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB 设置测试数据库
func SetupTestDB() *gorm.DB {
	cfg := GetTestConfig()
	
	// 构建数据库连接字符串
	dsn := cfg.Database.GetDSN()
	
	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 测试时使用静默模式
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(
		&models.User{},
		&models.Drama{},
		&models.Episode{},
		&models.Admin{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB(db *gorm.DB) {
	// 删除所有测试数据
	db.Exec("DELETE FROM episodes")
	db.Exec("DELETE FROM dramas")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM admins")
	
	// 重置自增ID
	db.Exec("ALTER TABLE episodes AUTO_INCREMENT = 1")
	db.Exec("ALTER TABLE dramas AUTO_INCREMENT = 1")
	db.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	db.Exec("ALTER TABLE admins AUTO_INCREMENT = 1")
}

// TruncateTable 清空指定表
func TruncateTable(db *gorm.DB, tableName string) {
	db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))
}

// BeginTransaction 开始事务
func BeginTransaction(db *gorm.DB) *gorm.DB {
	return db.Begin()
}

// RollbackTransaction 回滚事务
func RollbackTransaction(tx *gorm.DB) {
	tx.Rollback()
}

// CommitTransaction 提交事务
func CommitTransaction(tx *gorm.DB) {
	tx.Commit()
}

// WithTransaction 在事务中执行函数
func WithTransaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}