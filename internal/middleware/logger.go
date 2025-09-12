package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerConfig 日志配置
type LoggerConfig struct {
	// SkipPaths 跳过日志记录的路径
	SkipPaths []string
	// LogRequestBody 是否记录请求体
	LogRequestBody bool
	// LogResponseBody 是否记录响应体
	LogResponseBody bool
	// MaxBodySize 最大记录的请求/响应体大小
	MaxBodySize int64
}

// DefaultLoggerConfig 默认日志配置
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		SkipPaths: []string{
			"/health",
			"/ready",
			"/live",
			"/metrics",
		},
		LogRequestBody:  true,
		LogResponseBody: false,
		MaxBodySize:     1024 * 10, // 10KB
	}
}

// Logger 请求日志记录中间件
func Logger(config ...LoggerConfig) gin.HandlerFunc {
	conf := DefaultLoggerConfig()
	if len(config) > 0 {
		conf = config[0]
	}

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 检查是否跳过此路径
		for _, skipPath := range conf.SkipPaths {
			if param.Path == skipPath {
				return ""
			}
		}

		// 构建日志信息
		logData := map[string]interface{}{
			"timestamp":    param.TimeStamp.Format(time.RFC3339),
			"status":       param.StatusCode,
			"latency":      param.Latency.String(),
			"client_ip":    param.ClientIP,
			"method":       param.Method,
			"path":         param.Path,
			"user_agent":   param.Request.UserAgent(),
			"error":        param.ErrorMessage,
			"body_size":    param.BodySize,
		}

		// 添加请求 ID（如果存在）
		if requestID := param.Request.Header.Get("X-Request-ID"); requestID != "" {
			logData["request_id"] = requestID
		}

		// 添加用户信息（如果存在）
		if userID := param.Keys["user_id"]; userID != nil {
			logData["user_id"] = userID
		}
		if username := param.Keys["username"]; username != nil {
			logData["username"] = username
		}

		// 序列化为 JSON
		jsonData, _ := json.Marshal(logData)
		return string(jsonData) + "\n"
	})
}

// DetailedLogger 详细日志记录中间件（包含请求/响应体）
func DetailedLogger(config ...LoggerConfig) gin.HandlerFunc {
	conf := DefaultLoggerConfig()
	if len(config) > 0 {
		conf = config[0]
	}

	return func(c *gin.Context) {
		// 检查是否跳过此路径
		for _, skipPath := range conf.SkipPaths {
			if c.Request.URL.Path == skipPath {
				c.Next()
				return
			}
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 读取请求体
		var requestBody []byte
		if conf.LogRequestBody && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, conf.MaxBodySize))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建响应体写入器
		var responseBody *bytes.Buffer
		var writer gin.ResponseWriter = c.Writer
		if conf.LogResponseBody {
			responseBody = &bytes.Buffer{}
			writer = &responseBodyWriter{
				ResponseWriter: c.Writer,
				body:          responseBody,
			}
			c.Writer = writer
		}

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)

		// 构建完整路径
		if raw != "" {
			path = path + "?" + raw
		}

		// 构建日志数据
		logData := map[string]interface{}{
			"timestamp":   start.Format(time.RFC3339),
			"status":      c.Writer.Status(),
			"latency":     latency.String(),
			"latency_ms":  float64(latency.Nanoseconds()) / 1000000,
			"client_ip":   c.ClientIP(),
			"method":      c.Request.Method,
			"path":        path,
			"user_agent":  c.Request.UserAgent(),
			"body_size":   c.Writer.Size(),
		}

		// 添加请求头信息
		headers := make(map[string]string)
		for key, values := range c.Request.Header {
			if len(values) > 0 {
				// 过滤敏感头信息
				if key == "Authorization" || key == "Cookie" {
					headers[key] = "[REDACTED]"
				} else {
					headers[key] = values[0]
				}
			}
		}
		logData["headers"] = headers

		// 添加请求体
		if conf.LogRequestBody && len(requestBody) > 0 {
			// 尝试解析为 JSON
			var jsonBody interface{}
			if json.Unmarshal(requestBody, &jsonBody) == nil {
				logData["request_body"] = jsonBody
			} else {
				logData["request_body"] = string(requestBody)
			}
		}

		// 添加响应体
		if conf.LogResponseBody && responseBody != nil && responseBody.Len() > 0 {
			responseData := responseBody.Bytes()
			if len(responseData) > int(conf.MaxBodySize) {
				responseData = responseData[:conf.MaxBodySize]
			}
			
			// 尝试解析为 JSON
			var jsonResponse interface{}
			if json.Unmarshal(responseData, &jsonResponse) == nil {
				logData["response_body"] = jsonResponse
			} else {
				logData["response_body"] = string(responseData)
			}
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			logData["errors"] = c.Errors.Errors()
		}

		// 添加用户信息
		if userID, exists := c.Get("user_id"); exists {
			logData["user_id"] = userID
		}
		if username, exists := c.Get("username"); exists {
			logData["username"] = username
		}
		if role, exists := c.Get("role"); exists {
			logData["role"] = role
		}

		// 输出日志
		jsonData, _ := json.Marshal(logData)
		fmt.Println(string(jsonData))
	}
}

// responseBodyWriter 响应体写入器
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// RequestIDMiddleware 请求 ID 中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID 生成请求 ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}