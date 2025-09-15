package repository

import (
	"gin-mysql-api/internal/models"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(offset, limit int) ([]models.User, int64, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
}

// DramaRepository 短剧数据访问接口
type DramaRepository interface {
	Create(drama *models.Drama) error
	GetByID(id uint) (*models.Drama, error)
	GetByIDWithEpisodes(id uint) (*models.Drama, error)
	GetList(offset, limit int, genre string) ([]models.Drama, int64, error)
	Update(drama *models.Drama) error
	Delete(id uint) error
	IncrementViewCount(id uint) error
	GetByGenre(genre string, offset, limit int) ([]models.Drama, int64, error)
	GetActiveList(offset, limit int) ([]models.Drama, int64, error)
}

// EpisodeRepository 剧集数据访问接口
type EpisodeRepository interface {
	Create(episode *models.Episode) error
	GetByID(id uint) (*models.Episode, error)
	GetByIDWithDrama(id uint) (*models.Episode, error)
	GetByDramaID(dramaID uint) ([]models.Episode, error)
	GetByDramaIDPaginated(dramaID uint, offset, limit int) ([]models.Episode, int64, error)
	GetList(offset, limit int) ([]models.Episode, int64, error)
	Update(episode *models.Episode) error
	Delete(id uint) error
	IncrementViewCount(id uint) error
	GetMaxEpisodeNum(dramaID uint) (int, error)
	ExistsByDramaIDAndEpisodeNum(dramaID uint, episodeNum int) (bool, error)
}

// AdminRepository 管理员数据访问接口
type AdminRepository interface {
	Create(admin *models.Admin) error
	GetByID(id uint) (*models.Admin, error)
	GetByEmail(email string) (*models.Admin, error)
	GetByUsername(username string) (*models.Admin, error)
	Update(admin *models.Admin) error
	Delete(id uint) error
	List(offset, limit int) ([]models.Admin, int64, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
}
