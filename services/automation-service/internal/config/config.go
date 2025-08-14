package config

import (
	"os"
)

type Config struct {
	Server  ServerConfig  `json:"server"`
	Logging LoggingConfig `json:"logging"`
}

type ServerConfig struct {
	Port string `json:"port"`
	Mode string `json:"mode"`
}

type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8095"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}