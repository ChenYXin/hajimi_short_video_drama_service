package repository

import (
	"testing"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// DramaRepositoryTestSuite 短剧仓库测试套件
type DramaRepositoryTestSuite struct {
	suite.Suite
	db       *gorm.DB
	repo     DramaRepository
	factory  *testutil.Factory
}

// SetupSuite 设置测试套件
func (suite *DramaRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.SetupTestDB()
	suite.repo = NewDramaRepository(suite.db)
	suite.factory = testutil.NewFactory()
}

// TearDownSuite 清理测试套件
func (suite *DramaRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// SetupTest 每个测试前的设置
func (suite *DramaRepositoryTestSuite) SetupTest() {
	testutil.CleanupTestDB(suite.db)
}

// TestCreate 测试创建短剧
func (suite *DramaRepositoryTestSuite) TestCreate() {
	drama := suite.factory.Drama.CreateDrama()
	
	err := suite.repo.Create(drama)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), drama.ID)
	assert.NotZero(suite.T(), drama.CreatedAt)
	assert.NotZero(suite.T(), drama.UpdatedAt)
}

// TestGetByID 测试根据ID获取短剧
func (suite *DramaRepositoryTestSuite) TestGetByID() {
	// 创建测试短剧
	drama := suite.factory.Drama.CreateDrama()
	err := suite.repo.Create(drama)
	assert.NoError(suite.T(), err)

	// 获取短剧
	foundDrama, err := suite.repo.GetByID(drama.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), drama.Title, foundDrama.Title)
	assert.Equal(suite.T(), drama.Genre, foundDrama.Genre)

	// 测试不存在的短剧
	_, err = suite.repo.GetByID(999)
	assert.Error(suite.T(), err)
}

// TestGetByIDWithEpisodes 测试获取短剧及其剧集
func (suite *DramaRepositoryTestSuite) TestGetByIDWithEpisodes() {
	// 创建测试短剧
	drama := suite.factory.Drama.CreateDrama()
	err := suite.repo.Create(drama)
	assert.NoError(suite.T(), err)

	// 创建剧集
	episodeRepo := NewEpisodeRepository(suite.db)
	episodes := suite.factory.Episode.CreateEpisodes(drama.ID, 3)
	for _, episode := range episodes {
		err := episodeRepo.Create(episode)
		assert.NoError(suite.T(), err)
	}

	// 获取短剧及其剧集
	foundDrama, err := suite.repo.GetByIDWithEpisodes(drama.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), drama.Title, foundDrama.Title)
	assert.Len(suite.T(), foundDrama.Episodes, 3)
}

// TestGetList 测试获取短剧列表
func (suite *DramaRepositoryTestSuite) TestGetList() {
	// 创建多个测试短剧
	dramas := suite.factory.Drama.CreateDramas(5)
	for _, drama := range dramas {
		err := suite.repo.Create(drama)
		assert.NoError(suite.T(), err)
	}

	// 测试分页获取
	dramaList, total, err := suite.repo.GetList(0, 3, "")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(5), total)
	assert.Len(suite.T(), dramaList, 3)

	// 测试第二页
	dramaList, total, err = suite.repo.GetList(3, 3, "")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(5), total)
	assert.Len(suite.T(), dramaList, 2)
}

// TestGetByGenre 测试根据类型获取短剧
func (suite *DramaRepositoryTestSuite) TestGetByGenre() {
	// 创建不同类型的短剧
	comedyDrama := suite.factory.Drama.CreateDrama(func(d *models.Drama) {
		d.Genre = "喜剧"
		d.Title = "喜剧短剧"
	})
	actionDrama := suite.factory.Drama.CreateDrama(func(d *models.Drama) {
		d.Genre = "动作"
		d.Title = "动作短剧"
	})

	err := suite.repo.Create(comedyDrama)
	assert.NoError(suite.T(), err)
	err = suite.repo.Create(actionDrama)
	assert.NoError(suite.T(), err)

	// 测试获取喜剧类型
	comedyList, total, err := suite.repo.GetByGenre("喜剧", 0, 10)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), comedyList, 1)
	assert.Equal(suite.T(), "喜剧短剧", comedyList[0].Title)

	// 测试获取动作类型
	actionList, total, err := suite.repo.GetByGenre("动作", 0, 10)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), actionList, 1)
	assert.Equal(suite.T(), "动作短剧", actionList[0].Title)
}

// TestGetActiveList 测试获取激活状态的短剧列表
func (suite *DramaRepositoryTestSuite) TestGetActiveList() {
	// 创建不同状态的短剧
	activeDrama := suite.factory.Drama.CreateDrama(func(d *models.Drama) {
		d.Status = "active"
		d.Title = "激活短剧"
	})
	inactiveDrama := suite.factory.Drama.CreateDrama(func(d *models.Drama) {
		d.Status = "inactive"
		d.Title = "禁用短剧"
	})

	err := suite.repo.Create(activeDrama)
	assert.NoError(suite.T(), err)
	err = suite.repo.Create(inactiveDrama)
	assert.NoError(suite.T(), err)

	// 测试获取激活状态的短剧
	activeList, total, err := suite.repo.GetActiveList(0, 10)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), activeList, 1)
	assert.Equal(suite.T(), "激活短剧", activeList[0].Title)
}

// TestUpdate 测试更新短剧
func (suite *DramaRepositoryTestSuite) TestUpdate() {
	// 创建测试短剧
	drama := suite.factory.Drama.CreateDrama()
	err := suite.repo.Create(drama)
	assert.NoError(suite.T(), err)

	// 更新短剧信息
	drama.Title = "更新后的标题"
	drama.Description = "更新后的描述"
	err = suite.repo.Update(drama)
	assert.NoError(suite.T(), err)

	// 验证更新
	updatedDrama, err := suite.repo.GetByID(drama.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "更新后的标题", updatedDrama.Title)
	assert.Equal(suite.T(), "更新后的描述", updatedDrama.Description)
}

// TestDelete 测试删除短剧
func (suite *DramaRepositoryTestSuite) TestDelete() {
	// 创建测试短剧
	drama := suite.factory.Drama.CreateDrama()
	err := suite.repo.Create(drama)
	assert.NoError(suite.T(), err)

	// 删除短剧
	err = suite.repo.Delete(drama.ID)
	assert.NoError(suite.T(), err)

	// 验证删除
	_, err = suite.repo.GetByID(drama.ID)
	assert.Error(suite.T(), err)
}

// TestIncrementViewCount 测试增加观看次数
func (suite *DramaRepositoryTestSuite) TestIncrementViewCount() {
	// 创建测试短剧
	drama := suite.factory.Drama.CreateDrama()
	err := suite.repo.Create(drama)
	assert.NoError(suite.T(), err)

	initialViewCount := drama.ViewCount

	// 增加观看次数
	err = suite.repo.IncrementViewCount(drama.ID)
	assert.NoError(suite.T(), err)

	// 验证观看次数增加
	updatedDrama, err := suite.repo.GetByID(drama.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), initialViewCount+1, updatedDrama.ViewCount)
}

// TestDramaRepositoryTestSuite 运行短剧仓库测试套件
func TestDramaRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DramaRepositoryTestSuite))
}