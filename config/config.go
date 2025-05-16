package config

import (
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	JWT struct {
		Secret string
	}
	Server struct {
		Port string
	}
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("environment", "development")
	v.SetDefault("server.port", "8080")
	v.SetDefault("db.port", "5432")
	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.user", "go_auth")
	v.SetDefault("db.password", "go_auth")
	v.SetDefault("db.name", "go_auth")
	v.SetDefault("jwt.secret", "h13dpx8nFiWwLbhHuOEBLWhA6kfYwoP9UNU5MQlgoZQ0")
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	// Read from YAML file first (lower priority)
	v.SetConfigName("development")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.ErrLoadConfig
		}
	}

	// Read from .env file (higher priority)
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath("./config")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.ErrLoadConfig
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, errors.ErrLoadConfig
	}

	return &config, nil
}
