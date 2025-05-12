package repository

import (
	"context"
	"time"

	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/redis/go-redis/v9"
)

func (r *RedisRepository) AddToken(ctx context.Context, userID, token string, expiration time.Duration) error {
	err := r.client.Set(ctx, userID, token, expiration).Err()
	if err != nil {
		r.logger.Error("Error adding token",
			ports.F("error", err),
			ports.F("userID", userID),
			ports.F("token", token),
			ports.F("expiration", expiration),
		)
		return errors.ErrAddToken
	}
	return nil
}

func (r *RedisRepository) RemoveToken(ctx context.Context, userID string) error {
	err := r.client.Del(ctx, userID).Err()
	if err != nil {
		r.logger.Error("Error removing token",
			ports.F("error", err),
			ports.F("userID", userID),
		)
		return errors.ErrRemoveToken
	}
	return nil
}

func (r *RedisRepository) FindToken(ctx context.Context, userID string) (string, error) {
	val, err := r.client.Get(ctx, userID).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Error("Token not found",
				ports.F("error", err),
				ports.F("userID", userID),
			)
			return "", errors.ErrTokenNotFound
		}
		r.logger.Error("Error getting token",
			ports.F("error", err),
			ports.F("userID", userID),
		)
		return "", errors.ErrGetToken
	}
	return val, nil
}
