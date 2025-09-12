package repository

import (
	"testing"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// EpisodeRepositoryTestSuite 剧集仓库测试套件
type EpisodeRepositoryTestSuite struct {
	suite.Suite
	db          *gorm.DB
	repo        EpisodeRepository
	dramaRepo   DramaRepository
	factory     *testutil.Factory
	testDrama   *models.Drama
}

// SetupSuite 设置测试套件
func (suite *EpisodeRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.SetupTestDB()
	suite.repo = NewEpisodeRepository(suite.db)
	suite.dramaRepo = NewDramaRepository(suite.db)
	suite.factory = testutil.NewFactory()
}

// TearDownSuite 清理测试套件
func (suite *EpisodeRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// SetupTest 每个测试前的设置
func (suite *EpisodeRepositoryTestSuite) SetupTest() {
	testutil.CleanupTestDB(suite.db)
	
	// 创建测试短剧
	suite.testDrama = suite.factory.Drama.CreateDrama()
	err := suite.dramaRepo.Create(suite.testDrama)
	assert.NoError(suite.T(), err)
}

// TestCreate 测试创建剧集
func (suite *EpisodeRepositoryTestSuite) TestCreate() {
	episode := suite.factory.Episode.CreateEpisode(suite.testDrama.ID)
	
	err := suite.repo.Create(episode)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), episode.ID)
	assert.NotZero(suite.T(), episode.CreatedAt)
	assert.NotZero(suite.T(), episode.UpdatedAt)
}

// TestGetByID 测试根据ID获取剧集
func (suite *EpisodeRepositoryTestSuite) TestGetByID() {
	// 创建测试剧集
	episode := suite.factory.Episode.CreateEpisode(suite.testDrama.ID)
	err := suite.repo.Create(episode)
	assert.NoError(suite.T(), err)

	// 获取剧集
	foundEpisode, err := suite.repo.GetByID(episode.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), episode.Title, foundEpisode.Title)
	assert.Equal(suite.T(), episode.DramaID, foundEpisode.DramaID)

	// 测试不存在的剧集
	foundEpisode, err = suite.repo.GetByID(999)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundEpisode)
}

// TestGetByIDWithDrama 测试根据ID获取剧集（包含短剧信息）
func (suite *EpisodeRepositoryTestSuite) TestGetByIDWithDrama() {
	// 创建测试剧集
	episode := suite.factory.Episode.CreateEpisode(suite.testDrama.ID)
	err := suite.repo.Create(episode)
	assert.NoError(suite.T(), err)

	// 获取剧集（包含短剧信息）
	foundEpisode, err := suite.repo.GetByIDWithDrama(episode.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), episode.Title, foundEpisode.Title)
	assert.Equal(suite.T(), episode.DramaID, foundEpisode.DramaID)
	// 注意：这里可能需要根据实际的模型关联来调整断言
}

// TestGetByDramaID 测试根据短剧ID获取所有剧集
func (suite *EpisodeRepositoryTestSuite) TestGetByDramaID() {
	// 创建多个测试剧集
	episodes := suite.factory.Episode.CreateEpisodes(suite.testDrama.ID, 3)
	for _, episode := range episodes {
		err := suite.repo.Create(episode)
		assert.NoError(suite.T(), err)
	}

	// 获取剧集列表
	foundEpisodes, err := suite.repo.GetByDramaID(suite.testDrama.ID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), foundEpisodes, 3)

	// 验证按剧集号排序
	for i := 0; i < len(foundEpisodes)-1; i++ {
		assert.LessOrEqual(suite.T(), foundEpisodes[i].EpisodeNum, foundEpisodes[i+1].EpisodeNum)
	}
}

// TestGetByDramaIDPaginated 测试根据短剧ID获取剧集列表（分页）
func (suite *EpisodeRepositoryTestSuite) TestGetByDramaIDPaginated() {
	// 创建多个测试剧集
	episodes := suite.factory.Episode.CreateEpisodes(suite.testDrama.ID, 5)
	for _, episode := range episodes {
		err := suite.repo.Create(episode)
		assert.NoError(suite.T(), err)
	}

	// 测试分页获取
	episodeList, total, err := suite.repo.GetByDramaIDPaginated(suite.testDrama.ID, 0, 3)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(5), total)
	assert.Len(suite.T(), episodeList, 3)

	// 测试第二页
	episodeList, total, err = suite.repo.GetByDramaIDPaginated(suite.testDrama.ID, 3, 3)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(5), total)
	assert.Len(suite.T(), episodeList, 2)
}

// TestUpdate 测试更新剧集
func (suite *EpisodeRepositoryTestSuite) TestUpdate() {
	// 创建测试剧集
	episode := suite.factory.Episode.CreateEpisode(suite.testDrama.ID)
	err := suite.repo.Create(episode)
	assert.NoError(suite.T(), err)

	// 更新剧集信息
	episode.Title = "更新后的剧集标题"
	episode.Duration = 45
	err = suite.repo.Update(episode)
	assert.NoError(suite.T(), err)

	// 验证更新
	updatedEpisode, err := suite.repo.GetByID(episode.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "更新后的剧集标题", updatedEpisode.Title)
	assert.Equal(suite.T(), 45, updatedEpisode.Duration)
}

// TestDelete 测试删除剧集
func (suite *EpisodeRepositoryTestSuite) TestDelete() {
	// 创建测试剧集
	episode := suite.factory.Episode.CreateEpisode(suite.testDrama.ID)
	err := suite.repo.Create(episode)
	assert.NoError(suite.T(), err)

	// 删除剧集
	err = suite.repo.Delete(episode.ID)
	assert.NoError(suite.T(), err)

	// 验证删除
	foundEpisode, err := suite.repo.GetByID(episode.ID)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundEpisode)
}

// TestIncrementViewCount 测试增加观看次数
func (suite *EpisodeRepositoryTestSuite) TestIncrementViewCount() {
	// 创建测试剧集
	episode := suite.factory.Episode.CreateEpisode(suite.testDrama.ID)
	err := suite.repo.Create(episode)
	assert.NoError(suite.T(), err)

	initialViewCount := episode.ViewCount

	// 增加观看次数
	err = suite.repo.IncrementViewCount(episode.ID)
	assert.NoError(suite.T(), err)

	// 验证观看次数增加
	updatedEpisode, err := suite.repo.GetByID(episode.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), initialViewCount+1, updatedEpisode.ViewCount)
}

// TestGetMaxEpisodeNum 测试获取指定短剧的最大剧集号
func (suite *EpisodeRepositoryTestSuite) TestGetMaxEpisodeNum() {
	// 创建多个测试剧集
	episodes := suite.factory.Episode.CreateEpisodes(suite.testDrama.ID, 3)
	for _, episode := range episodes {
		err := suite.repo.Create(episode)
		assert.NoError(suite.T(), err)
	}

	// 获取最大剧集号
	maxEpisodeNum, err := suite.repo.GetMaxEpisodeNum(suite.testDrama.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, maxEpisodeNum)

	// 测试没有剧集的短剧
	anotherDrama := suite.factory.Drama.CreateDrama(func(d *models.Drama) {
		d.Title = "另一个短剧"
	})
	err = suite.dramaRepo.Create(anotherDrama)
	assert.NoError(suite.T(), err)

	maxEpisodeNum, err = suite.repo.GetMaxEpisodeNum(anotherDrama.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, maxEpisodeNum)
}

// TestExistsByDramaIDAndEpisodeNum 测试检查指定短剧的剧集号是否已存在
func (suite *EpisodeRepositoryTestSuite) TestExistsByDramaIDAndEpisodeNum() {
	// 创建测试剧集
	episode := suite.factory.Episode.CreateEpisode(suite.testDrama.ID, func(e *models.Episode) {
		e.EpisodeNum = 1
	})
	err := suite.repo.Create(episode)
	assert.NoError(suite.T(), err)

	// 测试存在的剧集号
	exists, err := suite.repo.ExistsByDramaIDAndEpisodeNum(suite.testDrama.ID, 1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// 测试不存在的剧集号
	exists, err = suite.repo.ExistsByDramaIDAndEpisodeNum(suite.testDrama.ID, 2)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)

	// 测试不存在的短剧ID
	exists, err = suite.repo.ExistsByDramaIDAndEpisodeNum(999, 1)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

// TestEpisodeRepositoryTestSuite 运行剧集仓库测试套件
func TestEpisodeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(EpisodeRepositoryTestSuite))
}