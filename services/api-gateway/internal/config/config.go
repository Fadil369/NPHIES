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
	}
	
	Kafka struct {
		Brokers []string
		Topics  KafkaTopics
	}
	
	RateLimit struct {
		RequestsPerMinute int
		BurstSize         int
	}
	
	JWT struct {
		Secret     string
		Expiration int
	}
	
	Auth struct {
		OAuthURL     string
		ClientID     string
		ClientSecret string
	}
	
	Services struct {
		EligibilityURL   string
		ClaimsURL        string
		TerminologyURL   string
	}
	
	Security struct {
		EnableMTLS     bool
		TLSCertPath    string
		TLSKeyPath     string
		TrustedCACerts []string
	}
	
	Monitoring struct {
		MetricsEnabled bool
		TracingEnabled bool
		LogLevel       string
	}
}

type KafkaTopics struct {
	ClaimsIntake        string
	EligibilityRequests string
	EligibilityResponses string
	PriorAuthRequests   string
	PriorAuthStatus     string
	FraudAlerts         string
	AuditTrail          string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Database configuration
	cfg.Database.URL = getEnv("POSTGRES_URL", "postgres://nphies:nphies_password@localhost:5432/nphies?sslmode=disable")
	cfg.Database.MaxConnections = getEnvInt("DB_MAX_CONNECTIONS", 10)
	cfg.Database.MaxIdleTime = getEnvInt("DB_MAX_IDLE_TIME", 30)
	cfg.Database.ConnectionRetry = getEnvInt("DB_CONNECTION_RETRY", 3)

	// Redis configuration
	cfg.Redis.URL = getEnv("REDIS_URL", "redis://localhost:6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.DB = getEnvInt("REDIS_DB", 0)

	// Kafka configuration
	cfg.Kafka.Brokers = []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
	cfg.Kafka.Topics = KafkaTopics{
		ClaimsIntake:         "claims.intake.v1",
		EligibilityRequests:  "eligibility.requests.v1",
		EligibilityResponses: "eligibility.responses.v1",
		PriorAuthRequests:    "priorauth.requests.v1",
		PriorAuthStatus:      "priorauth.status.v1",
		FraudAlerts:          "fraud.alerts.v1",
		AuditTrail:           "audit.trail.v1",
	}

	// Rate limiting
	cfg.RateLimit.RequestsPerMinute = getEnvInt("RATE_LIMIT_RPM", 500)
	cfg.RateLimit.BurstSize = getEnvInt("RATE_LIMIT_BURST", 100)

	// JWT configuration
	cfg.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	cfg.JWT.Expiration = getEnvInt("JWT_EXPIRATION", 3600) // 1 hour

	// OAuth configuration
	cfg.Auth.OAuthURL = getEnv("OAUTH_URL", "https://auth.nphies.sa")
	cfg.Auth.ClientID = getEnv("OAUTH_CLIENT_ID", "")
	cfg.Auth.ClientSecret = getEnv("OAUTH_CLIENT_SECRET", "")

	// Service URLs
	cfg.Services.EligibilityURL = getEnv("ELIGIBILITY_SERVICE_URL", "http://localhost:8090")
	cfg.Services.ClaimsURL = getEnv("CLAIMS_SERVICE_URL", "http://localhost:8092")
	cfg.Services.TerminologyURL = getEnv("TERMINOLOGY_SERVICE_URL", "http://localhost:8091")

	// Security configuration
	cfg.Security.EnableMTLS = getEnvBool("ENABLE_MTLS", false)
	cfg.Security.TLSCertPath = getEnv("TLS_CERT_PATH", "")
	cfg.Security.TLSKeyPath = getEnv("TLS_KEY_PATH", "")

	// Monitoring configuration
	cfg.Monitoring.MetricsEnabled = getEnvBool("METRICS_ENABLED", true)
	cfg.Monitoring.TracingEnabled = getEnvBool("TRACING_ENABLED", true)
	cfg.Monitoring.LogLevel = getEnv("LOG_LEVEL", "info")

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