package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	Environment string
	
	Database struct {
		URL             string
		MaxConnections  int
		MaxIdleTime     int
		ConnectionRetry int
	}
	
	Redis struct {
		URL      string
		Password string
		DB       int
		TTL      int // Cache TTL in seconds
	}
	
	Kafka struct {
		Brokers []string
		Topics  KafkaTopics
	}
	
	Monitoring struct {
		MetricsEnabled bool
		LogLevel       string
	}
	
	Business struct {
		CacheTTL         int // Cache TTL in seconds (5 minutes = 300)
		MaxResponseTime  int // Maximum response time in milliseconds
		EnableRuleEngine bool
	}
}

type KafkaTopics struct {
	EligibilityRequests  string
	EligibilityResponses string
	AuditTrail          string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:        getEnv("PORT", "8090"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Database configuration
	cfg.Database.URL = getEnv("POSTGRES_URL", "postgres://nphies:nphies_password@localhost:5432/eligibility?sslmode=disable")
	cfg.Database.MaxConnections = getEnvInt("DB_MAX_CONNECTIONS", 10)
	cfg.Database.MaxIdleTime = getEnvInt("DB_MAX_IDLE_TIME", 30)
	cfg.Database.ConnectionRetry = getEnvInt("DB_CONNECTION_RETRY", 3)

	// Redis configuration
	cfg.Redis.URL = getEnv("REDIS_URL", "redis://localhost:6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.DB = getEnvInt("REDIS_DB", 0)
	cfg.Redis.TTL = getEnvInt("REDIS_TTL", 300) // 5 minutes

	// Kafka configuration
	cfg.Kafka.Brokers = []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
	cfg.Kafka.Topics = KafkaTopics{
		EligibilityRequests:  "eligibility.requests.v1",
		EligibilityResponses: "eligibility.responses.v1",
		AuditTrail:           "audit.trail.v1",
	}

	// Monitoring configuration
	cfg.Monitoring.MetricsEnabled = getEnvBool("METRICS_ENABLED", true)
	cfg.Monitoring.LogLevel = getEnv("LOG_LEVEL", "info")

	// Business configuration
	cfg.Business.CacheTTL = getEnvInt("CACHE_TTL", 300)         // 5 minutes
	cfg.Business.MaxResponseTime = getEnvInt("MAX_RESPONSE_TIME", 900) // 900ms
	cfg.Business.EnableRuleEngine = getEnvBool("ENABLE_RULE_ENGINE", true)

	return cfg, nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}