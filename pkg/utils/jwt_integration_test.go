package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// JWTIntegrationTestSuite JWT 集成测试套件
type JWTIntegrationTestSuite struct {
	suite.Suite
	jwtManager *JWTManager
}

// SetupSuite 设置测试套件
func (suite *JWTIntegrationTestSuite) SetupSuite() {
	suite.jwtManager = NewJWTManager("integration-test-secret", time.Hour)
}

// TestJWTWorkflow 测试完整的 JWT 工作流程
func (suite *JWTIntegrationTestSuite) TestJWTWorkflow() {
	userID := uint(123)
	username := "integrationuser"
	role := "user"

	// 1. 生成 token
	token, err := suite.jwtManager.GenerateToken(userID, username, role)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)

	// 2. 验证 token
	claims, err := suite.jwtManager.VerifyToken(token)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, claims.UserID)
	assert.Equal(suite.T(), username, claims.Username)
	assert.Equal(suite.T(), role, claims.Role)

	// 3. 刷新 token
	newToken, err := suite.jwtManager.RefreshToken(token)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), newToken)

	// 4. 验证新 token
	newClaims, err := suite.jwtManager.VerifyToken(newToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, newClaims.UserID)
	assert.Equal(suite.T(), username, newClaims.Username)
	assert.Equal(suite.T(), role, newClaims.Role)

	// 5. 验证新 token 的过期时间更新了
	assert.True(suite.T(), newClaims.ExpiresAt.After(claims.ExpiresAt.Time))
}

// TestJWTWithDifferentRoles 测试不同角色的 JWT
func (suite *JWTIntegrationTestSuite) TestJWTWithDifferentRoles() {
	testCases := []struct {
		userID   uint
		username string
		role     string
	}{
		{1, "user1", "user"},
		{2, "admin1", "admin"},
		{3, "superadmin", "super_admin"},
	}

	for _, tc := range testCases {
		// 生成 token
		token, err := suite.jwtManager.GenerateToken(tc.userID, tc.username, tc.role)
		assert.NoError(suite.T(), err)

		// 验证 token
		claims, err := suite.jwtManager.VerifyToken(token)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), tc.userID, claims.UserID)
		assert.Equal(suite.T(), tc.username, claims.Username)
		assert.Equal(suite.T(), tc.role, claims.Role)
	}
}

// TestJWTExpiration 测试 JWT 过期
func (suite *JWTIntegrationTestSuite) TestJWTExpiration() {
	// 创建一个很短过期时间的 JWT 管理器
	shortJWTManager := NewJWTManager("test-secret", time.Millisecond*100)

	token, err := shortJWTManager.GenerateToken(1, "testuser", "user")
	assert.NoError(suite.T(), err)

	// 立即验证应该成功
	claims, err := shortJWTManager.VerifyToken(token)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(1), claims.UserID)

	// 等待 token 过期
	time.Sleep(time.Millisecond * 200)

	// 验证过期的 token 应该失败
	_, err = shortJWTManager.VerifyToken(token)
	assert.Error(suite.T(), err)
}

// TestJWTWithInvalidSecret 测试使用错误密钥验证 token
func (suite *JWTIntegrationTestSuite) TestJWTWithInvalidSecret() {
	// 使用正确密钥生成 token
	token, err := suite.jwtManager.GenerateToken(1, "testuser", "user")
	assert.NoError(suite.T(), err)

	// 使用错误密钥创建新的管理器
	wrongJWTManager := NewJWTManager("wrong-secret", time.Hour)

	// 使用错误密钥验证 token 应该失败
	_, err = wrongJWTManager.VerifyToken(token)
	assert.Error(suite.T(), err)
}

// TestPasswordHashingWorkflow 测试密码哈希工作流程
func (suite *JWTIntegrationTestSuite) TestPasswordHashingWorkflow() {
	password := "mySecurePassword123!"

	// 1. 哈希密码
	hashedPassword, err := HashPassword(password)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)
	assert.NotEqual(suite.T(), password, hashedPassword)

	// 2. 验证正确密码
	isValid := VerifyPassword(hashedPassword, password)
	assert.True(suite.T(), isValid)

	// 3. 验证错误密码
	isValid = VerifyPassword(hashedPassword, "wrongPassword")
	assert.False(suite.T(), isValid)

	// 4. 验证空密码
	isValid = VerifyPassword(hashedPassword, "")
	assert.False(suite.T(), isValid)
}

// TestPasswordHashingConsistency 测试密码哈希的一致性
func (suite *JWTIntegrationTestSuite) TestPasswordHashingConsistency() {
	password := "testPassword123"

	// 多次哈希同一个密码
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NotEqual(suite.T(), hash1, hash2) // bcrypt 每次生成的哈希都不同

	// 但都应该能验证原密码
	assert.True(suite.T(), VerifyPassword(hash1, password))
	assert.True(suite.T(), VerifyPassword(hash2, password))
}

// TestJWTIntegrationTestSuite 运行 JWT 集成测试套件
func TestJWTIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(JWTIntegrationTestSuite))
}