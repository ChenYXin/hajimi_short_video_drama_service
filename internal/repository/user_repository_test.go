package repository

import (
	"testing"

	"gin-mysql-api/internal/models"
	"gin-mysql-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// UserRepositoryTestSuite 用户仓库测试套件
type UserRepositoryTestSuite struct {
	suite.Suite
	db       *gorm.DB
	repo     UserRepository
	factory  *testutil.Factory
}

// SetupSuite 设置测试套件
func (suite *UserRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.SetupTestDB()
	suite.repo = NewUserRepository(suite.db)
	suite.factory = testutil.NewFactory()
}

// TearDownSuite 清理测试套件
func (suite *UserRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// SetupTest 每个测试前的设置
func (suite *UserRepositoryTestSuite) SetupTest() {
	testutil.CleanupTestDB(suite.db)
}

// TestCreate 测试创建用户
func (suite *UserRepositoryTestSuite) TestCreate() {
	user := suite.factory.User.CreateUser()
	
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), user.ID)
	assert.NotZero(suite.T(), user.CreatedAt)
	assert.NotZero(suite.T(), user.UpdatedAt)
}

// TestGetByID 测试根据ID获取用户
func (suite *UserRepositoryTestSuite) TestGetByID() {
	// 创建测试用户
	user := suite.factory.User.CreateUser()
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)

	// 获取用户
	foundUser, err := suite.repo.GetByID(user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Username, foundUser.Username)
	assert.Equal(suite.T(), user.Email, foundUser.Email)

	// 测试不存在的用户
	_, err = suite.repo.GetByID(999)
	assert.Error(suite.T(), err)
}

// TestGetByEmail 测试根据邮箱获取用户
func (suite *UserRepositoryTestSuite) TestGetByEmail() {
	// 创建测试用户
	user := suite.factory.User.CreateUser()
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)

	// 根据邮箱获取用户
	foundUser, err := suite.repo.GetByEmail(user.Email)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Username, foundUser.Username)
	assert.Equal(suite.T(), user.ID, foundUser.ID)

	// 测试不存在的邮箱
	_, err = suite.repo.GetByEmail("nonexistent@example.com")
	assert.Error(suite.T(), err)
}

// TestGetByUsername 测试根据用户名获取用户
func (suite *UserRepositoryTestSuite) TestGetByUsername() {
	// 创建测试用户
	user := suite.factory.User.CreateUser()
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)

	// 根据用户名获取用户
	foundUser, err := suite.repo.GetByUsername(user.Username)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Email, foundUser.Email)
	assert.Equal(suite.T(), user.ID, foundUser.ID)

	// 测试不存在的用户名
	_, err = suite.repo.GetByUsername("nonexistent")
	assert.Error(suite.T(), err)
}

// TestUpdate 测试更新用户
func (suite *UserRepositoryTestSuite) TestUpdate() {
	// 创建测试用户
	user := suite.factory.User.CreateUser()
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)

	// 更新用户信息
	user.Username = "updateduser"
	user.Phone = "09876543210"
	err = suite.repo.Update(user)
	assert.NoError(suite.T(), err)

	// 验证更新
	updatedUser, err := suite.repo.GetByID(user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "updateduser", updatedUser.Username)
	assert.Equal(suite.T(), "09876543210", updatedUser.Phone)
}

// TestDelete 测试删除用户
func (suite *UserRepositoryTestSuite) TestDelete() {
	// 创建测试用户
	user := suite.factory.User.CreateUser()
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)

	// 删除用户
	err = suite.repo.Delete(user.ID)
	assert.NoError(suite.T(), err)

	// 验证删除
	_, err = suite.repo.GetByID(user.ID)
	assert.Error(suite.T(), err)
}

// TestList 测试获取用户列表
func (suite *UserRepositoryTestSuite) TestList() {
	// 创建多个测试用户
	users := suite.factory.User.CreateUsers(5)
	for _, user := range users {
		err := suite.repo.Create(user)
		assert.NoError(suite.T(), err)
	}

	// 测试分页获取
	userList, total, err := suite.repo.List(0, 3)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(5), total)
	assert.Len(suite.T(), userList, 3)

	// 测试第二页
	userList, total, err = suite.repo.List(3, 3)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(5), total)
	assert.Len(suite.T(), userList, 2)
}

// TestExistsByEmail 测试邮箱是否存在
func (suite *UserRepositoryTestSuite) TestExistsByEmail() {
	// 创建测试用户
	user := suite.factory.User.CreateUser()
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)

	// 测试存在的邮箱
	exists, err := suite.repo.ExistsByEmail(user.Email)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// 测试不存在的邮箱
	exists, err = suite.repo.ExistsByEmail("nonexistent@example.com")
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

// TestExistsByUsername 测试用户名是否存在
func (suite *UserRepositoryTestSuite) TestExistsByUsername() {
	// 创建测试用户
	user := suite.factory.User.CreateUser()
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)

	// 测试存在的用户名
	exists, err := suite.repo.ExistsByUsername(user.Username)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// 测试不存在的用户名
	exists, err = suite.repo.ExistsByUsername("nonexistent")
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

// TestUserRepositoryTestSuite 运行用户仓库测试套件
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}