package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server    ServerConfig    `json:"server"`
	Database  DatabaseConfig  `json:"database"`
	Redis     RedisConfig     `json:"redis"`
	ML        MLConfig        `json:"ml"`
	Privacy   PrivacyConfig   `json:"privacy"`
	IoT       IoTConfig       `json:"iot"`
	Logging   LoggingConfig   `json:"logging"`
}

type ServerConfig struct {
	Port string `json:"port"`
	Mode string `json:"mode"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type MLConfig struct {
	ModelRegistry string  `json:"model_registry"`
	FeatureStore  string  `json:"feature_store"`
	BatchSize     int     `json:"batch_size"`
	ScoreThreshold float64 `json:"score_threshold"`
}

type PrivacyConfig struct {
	DifferentialPrivacy bool    `json:"differential_privacy"`
	Epsilon             float64 `json:"epsilon"`
	Delta               float64 `json:"delta"`
	AnonymizationLevel  string  `json:"anonymization_level"`
}

type IoTConfig struct {
	MQTTBroker   string `json:"mqtt_broker"`
	MQTTPort     int    `json:"mqtt_port"`
	MQTTUsername string `json:"mqtt_username"`
	MQTTPassword string `json:"mqtt_password"`
}

type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

func Load() (*Config, error) {
	batchSize, _ := strconv.Atoi(getEnv("ML_BATCH_SIZE", "100"))
	scoreThreshold, _ := strconv.ParseFloat(getEnv("ML_SCORE_THRESHOLD", "0.8"), 64)
	epsilon, _ := strconv.ParseFloat(getEnv("PRIVACY_EPSILON", "1.0"), 64)
	delta, _ := strconv.ParseFloat(getEnv("PRIVACY_DELTA", "1e-5"), 64)
	mqttPort, _ := strconv.Atoi(getEnv("IOT_MQTT_PORT", "1883"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8094"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "analytics_db"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       3,
		},
		ML: MLConfig{
			ModelRegistry:  getEnv("ML_MODEL_REGISTRY", "http://localhost:5000"),
			FeatureStore:   getEnv("ML_FEATURE_STORE", "http://localhost:6066"),
			BatchSize:      batchSize,
			ScoreThreshold: scoreThreshold,
		},
		Privacy: PrivacyConfig{
			DifferentialPrivacy: getEnv("PRIVACY_DIFFERENTIAL", "true") == "true",
			Epsilon:             epsilon,
			Delta:               delta,
			AnonymizationLevel:  getEnv("PRIVACY_ANONYMIZATION_LEVEL", "k-anonymity"),
		},
		IoT: IoTConfig{
			MQTTBroker:   getEnv("IOT_MQTT_BROKER", "localhost"),
			MQTTPort:     mqttPort,
			MQTTUsername: getEnv("IOT_MQTT_USERNAME", ""),
			MQTTPassword: getEnv("IOT_MQTT_PASSWORD", ""),
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