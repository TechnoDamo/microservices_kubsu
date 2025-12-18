package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	SessionTTL      time.Duration
	CleanupInterval time.Duration
	DatabaseURL     string
}

func Load() *Config {
	port := getEnv("PORT", "8080")
	sessionTTL := getEnvAsInt("SESSION_TTL", 24)          // hours
	cleanupInterval := getEnvAsInt("CLEANUP_INTERVAL", 1) // hours
	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:test@localhost:51507/kubsu_project_db?sslmode=disable")

	return &Config{
		Port:            port,
		SessionTTL:      time.Duration(sessionTTL) * time.Hour,
		CleanupInterval: time.Duration(cleanupInterval) * time.Hour,
		DatabaseURL:     databaseURL,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
