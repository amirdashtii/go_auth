package config

import (
	"os"
	"sync"

	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"environment"`
	DB struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"db"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
}

var (
	config *Config
	once   sync.Once
	configErr error
)

// LoadConfig returns a singleton instance of Config
func LoadConfig() (*Config, error) {
	once.Do(func() {
		config, configErr = loadConfig()
	})
	return config, configErr
}

// loadConfig is the internal function that actually loads the configuration
func loadConfig() (*Config, error) {
	v := viper.New()

	// Get environment from environment variable or default to development
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Set default values
	setDefaultValues(v)

	// Read from YAML file first (lower priority)
	if err := readYAMLConfig(v, env); err != nil {
		return nil, err
	}

	// Read from .env file (higher priority)
	if err := readEnvConfig(v); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, errors.ErrLoadConfig
	}

	// Validate environment matches
	if config.Environment != env {
		return nil, errors.ErrLoadConfig
	}

	return &config, nil
}

func setDefaultValues(v *viper.Viper) {
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
}

func readYAMLConfig(v *viper.Viper, env string) error {
	v.SetConfigName(env)
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return errors.ErrLoadConfig
		}
	}
	return nil
}

func readEnvConfig(v *viper.Viper) error {
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath("./config")
	
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return errors.ErrLoadConfig
		}
	}
	return nil
}
