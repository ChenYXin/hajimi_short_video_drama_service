package handler

import (
	"net/http"

	"gin-mysql-api/internal/service"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件处理器
type FileHandler struct {
	*BaseHandler
	fileService service.FileService
}

// NewFileHandler 创建文件处理器
func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{
		BaseHandler: NewBaseHandler(),
		fileService: fileService,
	}
}

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传文件（头像、封面、视频、缩略图等）
// @Tags 文件
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Param type formData string false "文件类型" Enums(avatar,cover,video,thumbnail) default(others)
// @Success 200 {object} models.APIResponse{data=models.FileUploadResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "请选择要上传的文件")
		return
	}
	defer file.Close()

	// 获取文件类型
	uploadType := c.PostForm("type")
	if uploadType == "" {
		uploadType = "others"
	}

	// 上传文件
	response, err := h.fileService.UploadFile(file, header, uploadType)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "文件上传成功", response)
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除已上传的文件
// @Tags 文件
// @Security BearerAuth
// @Produce json
// @Param path query string true "文件路径"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/upload [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		h.ErrorResponse(c, http.StatusBadRequest, "文件路径不能为空")
		return
	}

	err := h.fileService.DeleteFile(filePath)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponseWithMessage(c, "文件删除成功", nil)
}