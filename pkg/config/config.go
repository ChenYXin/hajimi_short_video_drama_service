package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host       string   `mapstructure:"host"`
	Port       int      `mapstructure:"port"`
	Mode       string   `mapstructure:"mode"`
	AllowedIPs []string `mapstructure:"allowedIPs"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parseTime"`
	Loc             string        `mapstructure:"loc"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	MaxSize      int      `mapstructure:"maxSize"`
	AllowedTypes []string `mapstructure:"allowedTypes"`
	UploadPath   string   `mapstructure:"uploadPath"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	Filename string `mapstructure:"filename"`
}

// LoadConfig 加载指定路径的配置文件
func LoadConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)

	// 设置环境变量前缀
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// 设置环境变量键名替换
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 转换时间单位
	config.Database.ConnMaxLifetime *= time.Second
	config.JWT.Expiration *= time.Hour

	return &config, nil
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		viper.AddConfigPath("./configs")
		viper.AddConfigPath(".")
	}

	// 设置环境变量前缀
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// 设置环境变量键名替换
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 转换时间单位
	config.Database.ConnMaxLifetime *= time.Second
	config.JWT.Expiration *= time.Hour

	return &config, nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)
}

// GetRedisAddr 获取Redis连接地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetServerAddr 获取服务器监听地址
func (c *ServerConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
