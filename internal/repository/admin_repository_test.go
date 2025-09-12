package repository

import (
	"testing"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// AdminRepositoryTestSuite 管理员仓库测试套件
type AdminRepositoryTestSuite struct {
	suite.Suite
	db       *gorm.DB
	repo     AdminRepository
	factory  *testutil.Factory
}

// SetupSuite 设置测试套件
func (suite *AdminRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.SetupTestDB()
	suite.repo = NewAdminRepository(suite.db)
	suite.factory = testutil.NewFactory()
}

// TearDownSuite 清理测试套件
func (suite *AdminRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// SetupTest 每个测试前的设置
func (suite *AdminRepositoryTestSuite) SetupTest() {
	testutil.CleanupTestDB(suite.db)
}

// TestCreate 测试创建管理员
func (suite *AdminRepositoryTestSuite) TestCreate() {
	admin := suite.factory.Admin.CreateAdmin()
	
	err := suite.repo.Create(admin)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), admin.ID)
	assert.NotZero(suite.T(), admin.CreatedAt)
	assert.NotZero(suite.T(), admin.UpdatedAt)
}

// TestGetByID 测试根据ID获取管理员
func (suite *AdminRepositoryTestSuite) TestGetByID() {
	// 创建测试管理员
	admin := suite.factory.Admin.CreateAdmin()
	err := suite.repo.Create(admin)
	assert.NoError(suite.T(), err)

	// 获取管理员
	foundAdmin, err := suite.repo.GetByID(admin.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), admin.Username, foundAdmin.Username)
	assert.Equal(suite.T(), admin.Email, foundAdmin.Email)

	// 测试不存在的管理员
	foundAdmin, err = suite.repo.GetByID(999)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundAdmin)
}

// TestGetByEmail 测试根据邮箱获取管理员
func (suite *AdminRepositoryTestSuite) TestGetByEmail() {
	// 创建测试管理员
	admin := suite.factory.Admin.CreateAdmin()
	err := suite.repo.Create(admin)
	assert.NoError(suite.T(), err)

	// 根据邮箱获取管理员
	foundAdmin, err := suite.repo.GetByEmail(admin.Email)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), admin.Username, foundAdmin.Username)
	assert.Equal(suite.T(), admin.ID, foundAdmin.ID)

	// 测试不存在的邮箱
	foundAdmin, err = suite.repo.GetByEmail("nonexistent@example.com")
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundAdmin)
}

// TestGetByUsername 测试根据用户名获取管理员
func (suite *AdminRepositoryTestSuite) TestGetByUsername() {
	// 创建测试管理员
	admin := suite.factory.Admin.CreateAdmin()
	err := suite.repo.Create(admin)
	assert.NoError(suite.T(), err)

	// 根据用户名获取管理员
	foundAdmin, err := suite.repo.GetByUsername(admin.Username)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), admin.Email, foundAdmin.Email)
	assert.Equal(suite.T(), admin.ID, foundAdmin.ID)

	// 测试不存在的用户名
	foundAdmin, err = suite.repo.GetByUsername("nonexistent")
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundAdmin)
}

// TestUpdate 测试更新管理员
func (suite *AdminRepositoryTestSuite) TestUpdate() {
	// 创建测试管理员
	admin := suite.factory.Admin.CreateAdmin()
	err := suite.repo.Create(admin)
	assert.NoError(suite.T(), err)

	// 更新管理员信息
	admin.Username = "updatedadmin"
	admin.Role = "superadmin"
	err = suite.repo.Update(admin)
	assert.NoError(suite.T(), err)

	// 验证更新
	updatedAdmin, err := suite.repo.GetByID(admin.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "updatedadmin", updatedAdmin.Username)
	assert.Equal(suite.T(), "superadmin", updatedAdmin.Role)
}

// TestDelete 测试删除管理员
func (suite *AdminRepositoryTestSuite) TestDelete() {
	// 创建测试管理员
	admin := suite.factory.Admin.CreateAdmin()
	err := suite.repo.Create(admin)
	assert.NoError(suite.T(), err)

	// 删除管理员
	err = suite.repo.Delete(admin.ID)
	assert.NoError(suite.T(), err)

	// 验证删除
	foundAdmin, err := suite.repo.GetByID(admin.ID)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundAdmin)
}

// TestList 测试获取管理员列表
func (suite *AdminRepositoryTestSuite) TestList() {
	// 创建多个测试管理员
	for i := 0; i < 5; i++ {
		admin := suite.factory.Admin.CreateAdmin(func(a *models.Admin) {
			a.Username = "testadmin" + string(rune(i+'0'))
			a.Email = "admin" + string(rune(i+'0')) + "@example.com"
		})
		err := suite.repo.Create(admin)
		assert.NoError(suite.T(), err)
	}

	// 测试分页获取
	adminList, total, err := suite.repo.List(0, 3)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(5), total)
	assert.Len(suite.T(), adminList, 3)

	// 测试第二页
	adminList, total, err = suite.repo.List(3, 3)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(5), total)
	assert.Len(suite.T(), adminList, 2)
}

// TestExistsByEmail 测试邮箱是否存在
func (suite *AdminRepositoryTestSuite) TestExistsByEmail() {
	// 创建测试管理员
	admin := suite.factory.Admin.CreateAdmin()
	err := suite.repo.Create(admin)
	assert.NoError(suite.T(), err)

	// 测试存在的邮箱
	exists, err := suite.repo.ExistsByEmail(admin.Email)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// 测试不存在的邮箱
	exists, err = suite.repo.ExistsByEmail("nonexistent@example.com")
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

// TestExistsByUsername 测试用户名是否存在
func (suite *AdminRepositoryTestSuite) TestExistsByUsername() {
	// 创建测试管理员
	admin := suite.factory.Admin.CreateAdmin()
	err := suite.repo.Create(admin)
	assert.NoError(suite.T(), err)

	// 测试存在的用户名
	exists, err := suite.repo.ExistsByUsername(admin.Username)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// 测试不存在的用户名
	exists, err := suite.repo.ExistsByUsername("nonexistent")
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

// TestAdminRepositoryTestSuite 运行管理员仓库测试套件
func TestAdminRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AdminRepositoryTestSuite))
}