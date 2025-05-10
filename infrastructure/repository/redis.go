package repository

import (
	"context"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	logger ports.Logger
}

func NewRedisRepository(logger ports.Logger) (*RedisRepository, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisRepository{client: client, logger: logger}, nil
}
