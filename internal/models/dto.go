package models

// 用户相关 DTO

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    string `json:"phone" validate:"omitempty,len=11"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
	Phone    string `json:"phone" validate:"omitempty,len=11"`
	Avatar   string `json:"avatar" validate:"omitempty"`
}

// 短剧相关 DTO

// CreateDramaRequest 创建短剧请求
type CreateDramaRequest struct {
	Title       string `json:"title" validate:"required,max=200"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
	Director    string `json:"director" validate:"max=100"`
	Actors      string `json:"actors" validate:"max=500"`
	Genre       string `json:"genre" validate:"required,max=100"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive draft"`
}

// UpdateDramaRequest 更新短剧请求
type UpdateDramaRequest struct {
	Title       string `json:"title" validate:"omitempty,max=200"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
	Director    string `json:"director" validate:"omitempty,max=100"`
	Actors      string `json:"actors" validate:"omitempty,max=500"`
	Genre       string `json:"genre" validate:"omitempty,max=100"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive draft"`
}

// 剧集相关 DTO

// CreateEpisodeRequest 创建剧集请求
type CreateEpisodeRequest struct {
	DramaID     uint   `json:"drama_id" validate:"required"`
	Title       string `json:"title" validate:"required,max=200"`
	EpisodeNum  int    `json:"episode_num" validate:"required,min=1"`
	Duration    int    `json:"duration" validate:"required,min=1"`
	VideoURL    string `json:"video_url"`
	Thumbnail   string `json:"thumbnail"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive draft"`
}

// UpdateEpisodeRequest 更新剧集请求
type UpdateEpisodeRequest struct {
	Title      string `json:"title" validate:"omitempty,max=200"`
	EpisodeNum int    `json:"episode_num" validate:"omitempty,min=1"`
	Duration   int    `json:"duration" validate:"omitempty,min=1"`
	VideoURL   string `json:"video_url"`
	Thumbnail  string `json:"thumbnail"`
	Status     string `json:"status" validate:"omitempty,oneof=active inactive draft"`
}

// 管理员相关 DTO

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// CreateAdminRequest 创建管理员请求
type CreateAdminRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=admin super_admin"`
}

// 通用响应 DTO

// APIResponse 通用 API 响应格式
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string      `json:"token"`
	ExpiresAt int64       `json:"expires_at"`
	User      interface{} `json:"user"`
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// 分页相关 DTO

// PaginatedUsers 分页用户响应
type PaginatedUsers struct {
	Users       []User `json:"users"`
	Total       int64  `json:"total"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
	TotalPages  int    `json:"total_pages"`
	HasNext     bool   `json:"has_next"`
	HasPrevious bool   `json:"has_previous"`
}

// PaginatedDramas 分页短剧响应
type PaginatedDramas struct {
	Dramas      []Drama `json:"dramas"`
	Total       int64   `json:"total"`
	Page        int     `json:"page"`
	PageSize    int     `json:"page_size"`
	TotalPages  int     `json:"total_pages"`
	HasNext     bool    `json:"has_next"`
	HasPrevious bool    `json:"has_previous"`
}

// PaginatedEpisodes 分页剧集响应
type PaginatedEpisodes struct {
	Episodes    []Episode `json:"episodes"`
	Total       int64     `json:"total"`
	Page        int       `json:"page"`
	PageSize    int       `json:"page_size"`
	TotalPages  int       `json:"total_pages"`
	HasNext     bool      `json:"has_next"`
	HasPrevious bool      `json:"has_previous"`
}

// PaginatedAdmins 分页管理员响应
type PaginatedAdmins struct {
	Admins      []Admin `json:"admins"`
	Total       int64   `json:"total"`
	Page        int     `json:"page"`
	PageSize    int     `json:"page_size"`
	TotalPages  int     `json:"total_pages"`
	HasNext     bool    `json:"has_next"`
	HasPrevious bool    `json:"has_previous"`
}

// JWTClaims JWT 声明结构（从 utils 包导入）
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}