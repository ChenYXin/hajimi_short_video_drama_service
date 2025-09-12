package service

import (
	"errors"
	"fmt"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/repository"
	"gin-mysql-api/pkg/utils"
)

// AdminService 管理服务接口
type AdminService interface {
	Login(req models.AdminLoginRequest) (*models.LoginResponse, error)
	CreateDrama(req models.CreateDramaRequest) (*models.Drama, error)
	UpdateDrama(id uint, req models.UpdateDramaRequest) (*models.Drama, error)
	DeleteDrama(id uint) error
	CreateEpisode(req models.CreateEpisodeRequest) (*models.Episode, error)
	UpdateEpisode(id uint, req models.UpdateEpisodeRequest) (*models.Episode, error)
	DeleteEpisode(id uint) error
	GetDramaList(page, pageSize int) (*models.PaginatedDramas, error)
	GetEpisodeList(dramaID uint, page, pageSize int) (*models.PaginatedEpisodes, error)
	CreateAdmin(req models.CreateAdminRequest) (*models.Admin, error)
	GetAdminList(page, pageSize int) (*models.PaginatedAdmins, error)
}

// adminService 管理服务实现
type adminService struct {
	adminRepo    repository.AdminRepository
	dramaRepo    repository.DramaRepository
	episodeRepo  repository.EpisodeRepository
	jwtManager   *utils.JWTManager
	cacheService CacheService
}

// NewAdminService 创建新的管理服务
func NewAdminService(
	adminRepo repository.AdminRepository,
	dramaRepo repository.DramaRepository,
	episodeRepo repository.EpisodeRepository,
	jwtManager *utils.JWTManager,
	cacheService CacheService,
) AdminService {
	return &adminService{
		adminRepo:    adminRepo,
		dramaRepo:    dramaRepo,
		episodeRepo:  episodeRepo,
		jwtManager:   jwtManager,
		cacheService: cacheService,
	}
}

// Login 管理员登录
func (s *adminService) Login(req models.AdminLoginRequest) (*models.LoginResponse, error) {
	// 根据用户名查找管理员
	admin, err := s.adminRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查管理员是否激活
	if !admin.IsActive {
		return nil, errors.New("管理员账户已被禁用")
	}

	// 验证密码
	if !utils.VerifyPassword(admin.Password, req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 JWT token
	token, err := s.jwtManager.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		return nil, fmt.Errorf("令牌生成失败: %w", err)
	}

	// 清除密码字段
	admin.Password = ""

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: 0, // 将在控制器中设置
		User:      admin,
	}, nil
}

// CreateDrama 创建短剧
func (s *adminService) CreateDrama(req models.CreateDramaRequest) (*models.Drama, error) {
	drama := &models.Drama{
		Title:       req.Title,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Director:    req.Director,
		Actors:      req.Actors,
		Genre:       req.Genre,
		Status:      req.Status,
	}

	// 设置默认状态
	if drama.Status == "" {
		drama.Status = "active"
	}

	err := s.dramaRepo.Create(drama)
	if err != nil {
		return nil, fmt.Errorf("创建短剧失败: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.DeletePattern("dramas:*")
		s.cacheService.DeletePattern("popular_dramas:*")
	}

	return drama, nil
}

// UpdateDrama 更新短剧
func (s *adminService) UpdateDrama(id uint, req models.UpdateDramaRequest) (*models.Drama, error) {
	// 获取现有短剧
	drama, err := s.dramaRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("短剧不存在: %w", err)
	}

	// 更新字段
	if req.Title != "" {
		drama.Title = req.Title
	}
	if req.Description != "" {
		drama.Description = req.Description
	}
	if req.CoverImage != "" {
		drama.CoverImage = req.CoverImage
	}
	if req.Director != "" {
		drama.Director = req.Director
	}
	if req.Actors != "" {
		drama.Actors = req.Actors
	}
	if req.Genre != "" {
		drama.Genre = req.Genre
	}
	if req.Status != "" {
		drama.Status = req.Status
	}

	err = s.dramaRepo.Update(drama)
	if err != nil {
		return nil, fmt.Errorf("更新短剧失败: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.Delete(fmt.Sprintf("drama:%d", id))
		s.cacheService.Delete(fmt.Sprintf("drama_with_episodes:%d", id))
		s.cacheService.DeletePattern("dramas:*")
		s.cacheService.DeletePattern("popular_dramas:*")
	}

	return drama, nil
}

// DeleteDrama 删除短剧
func (s *adminService) DeleteDrama(id uint) error {
	// 检查短剧是否存在
	_, err := s.dramaRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("短剧不存在: %w", err)
	}

	err = s.dramaRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("删除短剧失败: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.Delete(fmt.Sprintf("drama:%d", id))
		s.cacheService.Delete(fmt.Sprintf("drama_with_episodes:%d", id))
		s.cacheService.DeletePattern("dramas:*")
		s.cacheService.DeletePattern("episodes:drama:*")
		s.cacheService.DeletePattern("popular_dramas:*")
	}

	return nil
}

// CreateEpisode 创建剧集
func (s *adminService) CreateEpisode(req models.CreateEpisodeRequest) (*models.Episode, error) {
	// 检查短剧是否存在
	_, err := s.dramaRepo.GetByID(req.DramaID)
	if err != nil {
		return nil, fmt.Errorf("短剧不存在: %w", err)
	}

	// 检查剧集编号是否已存在
	exists, err := s.episodeRepo.ExistsByDramaIDAndEpisodeNum(req.DramaID, req.EpisodeNum)
	if err != nil {
		return nil, fmt.Errorf("检查剧集编号失败: %w", err)
	}
	if exists {
		return nil, errors.New("该剧集编号已存在")
	}

	episode := &models.Episode{
		DramaID:     req.DramaID,
		Title:       req.Title,
		EpisodeNum:  req.EpisodeNum,
		Duration:    req.Duration,
		VideoURL:    req.VideoURL,
		Thumbnail:   req.Thumbnail,
		Status:      req.Status,
	}

	// 设置默认状态
	if episode.Status == "" {
		episode.Status = "active"
	}

	err = s.episodeRepo.Create(episode)
	if err != nil {
		return nil, fmt.Errorf("创建剧集失败: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.Delete(fmt.Sprintf("drama_with_episodes:%d", req.DramaID))
		s.cacheService.DeletePattern(fmt.Sprintf("episodes:drama:%d:*", req.DramaID))
	}

	return episode, nil
}

// UpdateEpisode 更新剧集
func (s *adminService) UpdateEpisode(id uint, req models.UpdateEpisodeRequest) (*models.Episode, error) {
	// 获取现有剧集
	episode, err := s.episodeRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("剧集不存在: %w", err)
	}

	// 如果要更新剧集编号，检查是否已存在
	if req.EpisodeNum != 0 && req.EpisodeNum != episode.EpisodeNum {
		exists, err := s.episodeRepo.ExistsByDramaIDAndEpisodeNum(episode.DramaID, req.EpisodeNum)
		if err != nil {
			return nil, fmt.Errorf("检查剧集编号失败: %w", err)
		}
		if exists {
			return nil, errors.New("该剧集编号已存在")
		}
		episode.EpisodeNum = req.EpisodeNum
	}

	// 更新其他字段
	if req.Title != "" {
		episode.Title = req.Title
	}
	if req.Duration != 0 {
		episode.Duration = req.Duration
	}
	if req.VideoURL != "" {
		episode.VideoURL = req.VideoURL
	}
	if req.Thumbnail != "" {
		episode.Thumbnail = req.Thumbnail
	}
	if req.Status != "" {
		episode.Status = req.Status
	}

	err = s.episodeRepo.Update(episode)
	if err != nil {
		return nil, fmt.Errorf("更新剧集失败: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.Delete(fmt.Sprintf("episode:%d", id))
		s.cacheService.Delete(fmt.Sprintf("drama_with_episodes:%d", episode.DramaID))
		s.cacheService.DeletePattern(fmt.Sprintf("episodes:drama:%d:*", episode.DramaID))
	}

	return episode, nil
}

// DeleteEpisode 删除剧集
func (s *adminService) DeleteEpisode(id uint) error {
	// 获取剧集信息
	episode, err := s.episodeRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("剧集不存在: %w", err)
	}

	err = s.episodeRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("删除剧集失败: %w", err)
	}

	// 清除相关缓存
	if s.cacheService != nil {
		s.cacheService.Delete(fmt.Sprintf("episode:%d", id))
		s.cacheService.Delete(fmt.Sprintf("drama_with_episodes:%d", episode.DramaID))
		s.cacheService.DeletePattern(fmt.Sprintf("episodes:drama:%d:*", episode.DramaID))
	}

	return nil
}

// GetDramaList 获取短剧列表（管理员视图）
func (s *adminService) GetDramaList(page, pageSize int) (*models.PaginatedDramas, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	dramas, total, err := s.dramaRepo.GetList(offset, pageSize, "")
	if err != nil {
		return nil, fmt.Errorf("获取短剧列表失败: %w", err)
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	return &models.PaginatedDramas{
		Dramas:      dramas,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}, nil
}

// GetEpisodeList 获取剧集列表（管理员视图）
func (s *adminService) GetEpisodeList(dramaID uint, page, pageSize int) (*models.PaginatedEpisodes, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 检查短剧是否存在
	_, err := s.dramaRepo.GetByID(dramaID)
	if err != nil {
		return nil, fmt.Errorf("短剧不存在: %w", err)
	}

	offset := (page - 1) * pageSize
	episodes, total, err := s.episodeRepo.GetByDramaIDPaginated(dramaID, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取剧集列表失败: %w", err)
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	return &models.PaginatedEpisodes{
		Episodes:    episodes,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}, nil
}

// CreateAdmin 创建管理员
func (s *adminService) CreateAdmin(req models.CreateAdminRequest) (*models.Admin, error) {
	// 检查用户名是否已存在
	exists, err := s.adminRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	exists, err = s.adminRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return nil, errors.New("邮箱已被注册")
	}

	// 对密码进行哈希处理
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码处理失败: %w", err)
	}

	admin := &models.Admin{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		IsActive: true,
	}

	// 设置默认角色
	if admin.Role == "" {
		admin.Role = "admin"
	}

	err = s.adminRepo.Create(admin)
	if err != nil {
		return nil, fmt.Errorf("创建管理员失败: %w", err)
	}

	// 清除密码字段
	admin.Password = ""
	return admin, nil
}

// GetAdminList 获取管理员列表
func (s *adminService) GetAdminList(page, pageSize int) (*models.PaginatedAdmins, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	admins, total, err := s.adminRepo.List(offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取管理员列表失败: %w", err)
	}

	// 清除所有管理员的密码字段
	for i := range admins {
		admins[i].Password = ""
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	return &models.PaginatedAdmins{
		Admins:      admins,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}, nil
}