package service

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestCacheService_Set(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cacheService := NewCacheService(db)

	t.Run("成功设置缓存", func(t *testing.T) {
		key := "test:key"
		value := "test value"
		expiration := time.Hour

		mock.ExpectSet(key, value, expiration).SetVal("OK")

		err := cacheService.Set(key, value, expiration)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCacheService_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cacheService := NewCacheService(db)

	t.Run("成功获取缓存", func(t *testing.T) {
		key := "test:key"
		expectedValue := "test value"

		mock.ExpectGet(key).SetVal(expectedValue)

		value, err := cacheService.Get(key)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("缓存不存在", func(t *testing.T) {
		key := "nonexistent:key"

		mock.ExpectGet(key).SetErr(redis.Nil)

		value, err := cacheService.Get(key)

		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
		assert.Empty(t, value)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCacheService_Delete(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cacheService := NewCacheService(db)

	t.Run("成功删除缓存", func(t *testing.T) {
		key := "test:key"

		mock.ExpectDel(key).SetVal(1)

		err := cacheService.Delete(key)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCacheService_Exists(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cacheService := NewCacheService(db)

	t.Run("缓存存在", func(t *testing.T) {
		key := "test:key"

		mock.ExpectExists(key).SetVal(1)

		exists, err := cacheService.Exists(key)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("缓存不存在", func(t *testing.T) {
		key := "nonexistent:key"

		mock.ExpectExists(key).SetVal(0)

		exists, err := cacheService.Exists(key)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCacheService_SetJSON(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cacheService := NewCacheService(db)

	t.Run("成功设置JSON缓存", func(t *testing.T) {
		key := "test:json"
		value := map[string]interface{}{
			"name": "测试",
			"age":  25,
		}
		expiration := time.Hour

		jsonData, _ := json.Marshal(value)
		mock.ExpectSet(key, jsonData, expiration).SetVal("OK")

		err := cacheService.SetJSON(key, value, expiration)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCacheService_GetJSON(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cacheService := NewCacheService(db)

	t.Run("成功获取JSON缓存", func(t *testing.T) {
		key := "test:json"
		originalValue := map[string]interface{}{
			"name": "测试",
			"age":  float64(25), // JSON 数字会被解析为 float64
		}

		jsonData, _ := json.Marshal(originalValue)
		mock.ExpectGet(key).SetVal(string(jsonData))

		var result map[string]interface{}
		err := cacheService.GetJSON(key, &result)

		assert.NoError(t, err)
		assert.Equal(t, originalValue, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCacheService_DeletePattern(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cacheService := NewCacheService(db)

	t.Run("成功删除匹配模式的缓存", func(t *testing.T) {
		pattern := "test:*"
		keys := []string{"test:key1", "test:key2", "test:key3"}

		mock.ExpectKeys(pattern).SetVal(keys)
		mock.ExpectDel(keys...).SetVal(int64(len(keys)))

		err := cacheService.DeletePattern(pattern)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("没有匹配的键", func(t *testing.T) {
		pattern := "nonexistent:*"
		keys := []string{}

		mock.ExpectKeys(pattern).SetVal(keys)

		err := cacheService.DeletePattern(pattern)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCacheService_Increment(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cacheService := NewCacheService(db)

	t.Run("成功递增计数器", func(t *testing.T) {
		key := "counter:test"
		expectedValue := int64(1)

		mock.ExpectIncr(key).SetVal(expectedValue)

		value, err := cacheService.Increment(key)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCacheService_Expire(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cacheService := NewCacheService(db)

	t.Run("成功设置过期时间", func(t *testing.T) {
		key := "test:key"
		expiration := time.Hour

		mock.ExpectExpire(key, expiration).SetVal(true)

		err := cacheService.Expire(key, expiration)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}