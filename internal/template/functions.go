package template

import (
	"html/template"
	"strings"
	"time"
)

// GetFuncMap 返回模板函数映射
func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		// 字符串处理函数
		"substr": func(s string, start, length int) string {
			if start < 0 || start >= len(s) {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		"trim":  strings.TrimSpace,
		
		// 数学函数
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		
		// 比较函数
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"ne": func(a, b interface{}) bool {
			return a != b
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
		"gte": func(a, b int) bool {
			return a >= b
		},
		"lte": func(a, b int) bool {
			return a <= b
		},
		
		// 时间格式化函数
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"formatTime": func(t time.Time) string {
			return t.Format("15:04:05")
		},
		"timeAgo": func(t time.Time) string {
			duration := time.Since(t)
			
			if duration < time.Minute {
				return "刚刚"
			} else if duration < time.Hour {
				return formatDuration(int(duration.Minutes())) + "分钟前"
			} else if duration < 24*time.Hour {
				return formatDuration(int(duration.Hours())) + "小时前"
			} else if duration < 30*24*time.Hour {
				return formatDuration(int(duration.Hours()/24)) + "天前"
			} else {
				return t.Format("2006-01-02")
			}
		},
		
		// 文件大小格式化
		"formatFileSize": func(bytes int64) string {
			const unit = 1024
			if bytes < unit {
				return formatDuration(int(bytes)) + " B"
			}
			div, exp := int64(unit), 0
			for n := bytes / unit; n >= unit; n /= unit {
				div *= unit
				exp++
			}
			return formatFloat(float64(bytes)/float64(div), 1) + " " + "KMGTPE"[exp:exp+1] + "B"
		},
		
		// 时长格式化（分钟转换为小时分钟）
		"formatDuration": func(minutes int) string {
			if minutes < 60 {
				return formatDuration(minutes) + "分钟"
			}
			hours := minutes / 60
			mins := minutes % 60
			if mins == 0 {
				return formatDuration(hours) + "小时"
			}
			return formatDuration(hours) + "小时" + formatDuration(mins) + "分钟"
		},
		
		// 数组/切片操作
		"len": func(v interface{}) int {
			switch val := v.(type) {
			case []interface{}:
				return len(val)
			case string:
				return len(val)
			default:
				return 0
			}
		},
		"first": func(v interface{}) interface{} {
			switch val := v.(type) {
			case []interface{}:
				if len(val) > 0 {
					return val[0]
				}
			}
			return nil
		},
		"last": func(v interface{}) interface{} {
			switch val := v.(type) {
			case []interface{}:
				if len(val) > 0 {
					return val[len(val)-1]
				}
			}
			return nil
		},
		
		// 条件函数
		"default": func(defaultValue, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		
		// HTML 安全函数
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeCSS": func(s string) template.CSS {
			return template.CSS(s)
		},
		"safeJS": func(s string) template.JS {
			return template.JS(s)
		},
		
		// URL 处理函数
		"urlQuery": func(key, value string, params map[string]string) string {
			if params == nil {
				params = make(map[string]string)
			}
			params[key] = value
			
			var parts []string
			for k, v := range params {
				if v != "" {
					parts = append(parts, k+"="+v)
				}
			}
			
			if len(parts) > 0 {
				return "?" + strings.Join(parts, "&")
			}
			return ""
		},
		
		// 状态相关函数
		"statusClass": func(status string) string {
			switch status {
			case "active":
				return "success"
			case "inactive":
				return "danger"
			case "draft":
				return "warning"
			default:
				return "secondary"
			}
		},
		"statusText": func(status string) string {
			switch status {
			case "active":
				return "激活"
			case "inactive":
				return "禁用"
			case "draft":
				return "草稿"
			default:
				return "未知"
			}
		},
		
		// 分页函数
		"pageRange": func(current, total int) []int {
			var pages []int
			start := current - 2
			if start < 1 {
				start = 1
			}
			end := current + 2
			if end > total {
				end = total
			}
			for i := start; i <= end; i++ {
				pages = append(pages, i)
			}
			return pages
		},
	}
}

// 辅助函数
func formatDuration(n int) string {
	if n < 10 {
		return "0" + string(rune(n+'0'))
	}
	return string(rune(n/10+'0')) + string(rune(n%10+'0'))
}

func formatFloat(f float64, precision int) string {
	format := "%." + string(rune(precision+'0')) + "f"
	return strings.TrimRight(strings.TrimRight(template.HTMLEscapeString(template.JSEscapeString(format)), "0"), ".")
}