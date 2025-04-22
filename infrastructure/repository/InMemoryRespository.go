package repository

import (
	"context"
	"time"
)

func (r *RedisRepository) AddToken(userID, token string, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, userID, token, expiration).Err()
}

func (r *RedisRepository) RemoveToken(userID string) error {
	ctx := context.Background()
	return r.client.Del(ctx, userID).Err()
}

func (r *RedisRepository) FindToken(userID string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, userID).Result()

}
