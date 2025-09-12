package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gin-mysql-api/internal/models"
)

// FileService 文件服务接口
type FileService interface {
	UploadFile(file multipart.File, header *multipart.FileHeader, uploadType string) (*models.FileUploadResponse, error)
	DeleteFile(filePath string) error
	GetFileURL(filePath string) string
	ValidateFileType(filename string, allowedTypes []string) bool
	ValidateFileSize(size int64, maxSize int64) bool
}

// fileService 文件服务实现
type fileService struct {
	uploadPath string
	baseURL    string
	maxSize    int64
	allowedTypes []string
}

// NewFileService 创建新的文件服务
func NewFileService(uploadPath, baseURL string, maxSize int64, allowedTypes []string) FileService {
	return &fileService{
		uploadPath:   uploadPath,
		baseURL:      baseURL,
		maxSize:      maxSize,
		allowedTypes: allowedTypes,
	}
}

// UploadFile 上传文件
func (s *fileService) UploadFile(file multipart.File, header *multipart.FileHeader, uploadType string) (*models.FileUploadResponse, error) {
	// 验证文件大小
	if !s.ValidateFileSize(header.Size, s.maxSize) {
		return nil, fmt.Errorf("文件大小超过限制，最大允许 %d MB", s.maxSize/(1024*1024))
	}

	// 验证文件类型
	if !s.ValidateFileType(header.Filename, s.allowedTypes) {
		return nil, fmt.Errorf("不支持的文件类型，允许的类型: %s", strings.Join(s.allowedTypes, ", "))
	}

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), generateRandomString(8), ext)

	// 根据上传类型创建子目录
	var subDir string
	switch uploadType {
	case "avatar":
		subDir = "avatars"
	case "cover":
		subDir = "covers"
	case "video":
		subDir = "videos"
	case "thumbnail":
		subDir = "thumbnails"
	default:
		subDir = "others"
	}

	// 创建完整的文件路径
	fullDir := filepath.Join(s.uploadPath, subDir)
	err := os.MkdirAll(fullDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	filePath := filepath.Join(fullDir, filename)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 生成文件 URL
	relativePath := filepath.Join(subDir, filename)
	fileURL := s.GetFileURL(relativePath)

	return &models.FileUploadResponse{
		URL:      fileURL,
		Filename: filename,
		Size:     header.Size,
	}, nil
}

// DeleteFile 删除文件
func (s *fileService) DeleteFile(filePath string) error {
	fullPath := filepath.Join(s.uploadPath, filePath)
	
	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // 文件不存在，认为删除成功
	}

	return os.Remove(fullPath)
}

// GetFileURL 获取文件 URL
func (s *fileService) GetFileURL(filePath string) string {
	return fmt.Sprintf("%s/uploads/%s", s.baseURL, strings.ReplaceAll(filePath, "\\", "/"))
}

// ValidateFileType 验证文件类型
func (s *fileService) ValidateFileType(filename string, allowedTypes []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != "" && ext[0] == '.' {
		ext = ext[1:] // 移除点号
	}

	for _, allowedType := range allowedTypes {
		if ext == strings.ToLower(allowedType) {
			return true
		}
	}

	return false
}

// ValidateFileSize 验证文件大小
func (s *fileService) ValidateFileSize(size int64, maxSize int64) bool {
	return size <= maxSize
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}