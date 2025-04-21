package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	DBHost           string
	DBPort           int
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	JWTRefreshSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           dbPort,
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "product_management"),
		JWTSecret:        getEnv("JWT_SECRET", "01964c7b_9461_735b_82af_c02f626b7066"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "01964c7b_9461_735b_82af_c02f626b7066SASS"),
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
