package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
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
	v.SetDefault("server.port", "8080")
	v.SetDefault("db.port", "5432")
	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.user", "go_auth")
	v.SetDefault("db.password", "go_auth")
	v.SetDefault("db.name", "go_auth")
	v.SetDefault("jwt.secret", "h13dpx8nFiWwLbhHuOEBLWhA6kfYwoP9UNU5MQlgoZQ0")
	v.SetDefault("redis.Addr", "localhost:6379")
	v.SetDefault("redis.Password", "")
	v.SetDefault("redis.DB", 0)
	

	// Read from YAML file
	v.SetConfigName("development")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read YAML config file: %w", err)
		}
	}

	// Read from .env file
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath("./config")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read .env file: %wformat", err)
		}
	} else {
		v.AutomaticEnv()
		v.Set("server.port", v.GetString("SERVER_PORT"))
		v.Set("db.host", v.GetString("DB_HOST"))
		v.Set("db.port", v.GetString("DB_PORT"))
		v.Set("db.user", v.GetString("DB_USER"))
		v.Set("db.password", v.GetString("DB_PASSWORD"))
		v.Set("db.name", v.GetString("DB_NAME"))
		v.Set("jwt.secret", v.GetString("JWT_SECRET"))
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
