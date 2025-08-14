package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server     ServerConfig     `json:"server"`
	Database   DatabaseConfig   `json:"database"`
	Redis      RedisConfig      `json:"redis"`
	Blockchain BlockchainConfig `json:"blockchain"`
	Logging    LoggingConfig    `json:"logging"`
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

type BlockchainConfig struct {
	Network    string `json:"network"`
	NodeURL    string `json:"node_url"`
	PrivateKey string `json:"private_key"`
	ChainID    int    `json:"chain_id"`
	GasLimit   int    `json:"gas_limit"`
}

type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

func Load() (*Config, error) {
	chainID, _ := strconv.Atoi(getEnv("BLOCKCHAIN_CHAIN_ID", "1337"))
	gasLimit, _ := strconv.Atoi(getEnv("BLOCKCHAIN_GAS_LIMIT", "3000000"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8093"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "wallet_db"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       2,
		},
		Blockchain: BlockchainConfig{
			Network:    getEnv("BLOCKCHAIN_NETWORK", "hyperledger"),
			NodeURL:    getEnv("BLOCKCHAIN_NODE_URL", "http://localhost:7051"),
			PrivateKey: getEnv("BLOCKCHAIN_PRIVATE_KEY", ""),
			ChainID:    chainID,
			GasLimit:   gasLimit,
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