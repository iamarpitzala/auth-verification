package config

import (
	"os"
)

type Config struct {
	DatabaseURL     string
	JWTSecret       string
	Port            string
	EmailWebhookURL string
}

func Load() *Config {
	return &Config{
		DatabaseURL:     getEnv("DATABASE_URL", "postgresql://postgres:password@localhost:5432/auth_db"),
		JWTSecret:       getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		Port:            getEnv("PORT", "8080"),
		EmailWebhookURL: getEnv("EMAIL_WEBHOOK_URL", "https://script.google.com/macros/s/AKfycbwvkIqfvkERu_ATLz4Ci4LtF7_VKTIZ7eSZog9A9kLAR7AGaJmcMkOjqXWgirFDvvw/exec"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
