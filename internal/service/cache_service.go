package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheService Redis 缓存服务接口
type CacheService interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	SetJSON(key string, value interface{}, expiration time.Duration) error
	GetJSON(key string, dest interface{}) error
	DeletePattern(pattern string) error
	Increment(key string) (int64, error)
	Expire(key string, expiration time.Duration) error
}

// cacheService Redis 缓存服务实现
type cacheService struct {
	client *redis.Client
	ctx    context.Context
}

// NewCacheService 创建新的缓存服务
func NewCacheService(client *redis.Client) CacheService {
	return &cacheService{
		client: client,
		ctx:    context.Background(),
	}
}

// Set 设置缓存值
func (s *cacheService) Set(key string, value interface{}, expiration time.Duration) error {
	return s.client.Set(s.ctx, key, value, expiration).Err()
}

// Get 获取缓存值
func (s *cacheService) Get(key string) (string, error) {
	return s.client.Get(s.ctx, key).Result()
}

// Delete 删除缓存
func (s *cacheService) Delete(key string) error {
	return s.client.Del(s.ctx, key).Err()
}

// Exists 检查缓存是否存在
func (s *cacheService) Exists(key string) (bool, error) {
	count, err := s.client.Exists(s.ctx, key).Result()
	return count > 0, err
}

// SetJSON 设置 JSON 格式的缓存
func (s *cacheService) SetJSON(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(s.ctx, key, jsonData, expiration).Err()
}

// GetJSON 获取 JSON 格式的缓存
func (s *cacheService) GetJSON(key string, dest interface{}) error {
	jsonData, err := s.client.Get(s.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonData), dest)
}

// DeletePattern 根据模式删除缓存
func (s *cacheService) DeletePattern(pattern string) error {
	keys, err := s.client.Keys(s.ctx, pattern).Result()
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return s.client.Del(s.ctx, keys...).Err()
	}
	
	return nil
}

// Increment 递增计数器
func (s *cacheService) Increment(key string) (int64, error) {
	return s.client.Incr(s.ctx, key).Result()
}

// Expire 设置过期时间
func (s *cacheService) Expire(key string, expiration time.Duration) error {
	return s.client.Expire(s.ctx, key, expiration).Err()
}