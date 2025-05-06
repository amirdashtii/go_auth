package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *RedisRepository) AddToken(userID, token string, expiration time.Duration) error {
	ctx := context.Background()
	err := r.client.Set(ctx, userID, token, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) RemoveToken(userID string) error {
	ctx := context.Background()
	err := r.client.Del(ctx, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) FindToken(userID string) (string, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, userID).Result()
	if err != nil {
		if err == redis.Nil {
			return "", err
		}
		return "", err
	}
	return val, nil
}
