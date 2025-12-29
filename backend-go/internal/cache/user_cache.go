package cache

import (
	"context"
	"quiz-game-backend/internal/dto"
	"time"
)

const (
	UserCachePrefix = "user:"
	UserCacheTTL    = 24 * time.Hour
)

type UserCache struct {
	redis *RedisClient
}

func NewUserCache(redis *RedisClient) *UserCache {
	return &UserCache{redis: redis}
}
func (c *UserCache) Get(ctx context.Context, userID string) (*dto.UserResponse, error) {
	var user dto.UserResponse
	err := c.redis.Get(ctx, UserCachePrefix+userID, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (c *UserCache) Set(ctx context.Context, user *dto.UserResponse) error {
	return c.redis.Set(ctx, UserCachePrefix+user.ID, user, UserCacheTTL)
}
func (c *UserCache) Delete(ctx context.Context, userID string) error {
	return c.redis.Delete(ctx, UserCachePrefix+userID)
}
