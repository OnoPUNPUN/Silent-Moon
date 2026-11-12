package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv         string
	HTTPPort       string
	DatabaseURL    string
	JWTAccessSecret string
	JWTAccessTTL   time.Duration
	JWTRefreshTTL  time.Duration
	R2AccountID    string
	R2AccessKeyID  string
	R2SecretKey    string
	R2BucketName   string
	R2Endpoint     string
	AllowedOrigins []string
}

func Load() *Config {
	return &Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/reading_library?sslmode=disable"),
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "super-secret-development-key-change-me"),
		JWTAccessTTL:    getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute),
		JWTRefreshTTL:   getEnvDuration("JWT_REFRESH_TTL", 30*24*time.Hour),
		R2AccountID:     getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:   getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretKey:     getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2BucketName:    getEnv("R2_BUCKET_NAME", "reading-library"),
		R2Endpoint:      getEnv("R2_ENDPOINT", ""),
		AllowedOrigins:  splitOrigins(getEnv("ALLOWED_ORIGINS", "*")),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && strings.TrimSpace(val) != "" {
		return val
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func splitOrigins(origins string) []string {
	if origins == "" || origins == "*" {
		return []string{"*"}
	}
	parts := strings.Split(origins, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}
