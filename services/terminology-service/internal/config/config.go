package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port     string
	LogLevel string
	
	// Database configuration
	DatabaseURL      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	
	// Redis configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	
	// Cache configuration
	CacheTTL       int // in seconds
	CacheKeyPrefix string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8091"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		
		// Database
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		DatabaseHost:     getEnv("DB_HOST", "localhost"),
		DatabasePort:     getEnv("DB_PORT", "5432"),
		DatabaseName:     getEnv("DB_NAME", "nphies_terminology"),
		DatabaseUser:     getEnv("DB_USER", "postgres"),
		DatabasePassword: getEnv("DB_PASSWORD", "password"),
		
		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
		
		// Cache
		CacheTTL:       3600, // 1 hour default
		CacheKeyPrefix: "nphies:terminology:",
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}