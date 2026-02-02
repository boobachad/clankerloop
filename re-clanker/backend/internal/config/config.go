package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	DatabaseURL      string
	AIProvider       string // "openrouter" or "gemini"
	OpenRouterAPIKey string
	GeminiAPIKey     string
	Port             string
	CORSOrigins      string
	LogLevel         string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		AIProvider:       getEnvOrDefault("AI_PROVIDER", "openrouter"),
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
		Port:             getEnvOrDefault("PORT", "8080"),
		CORSOrigins:      getEnvOrDefault("CORS_ORIGINS", "http://localhost:3000"),
		LogLevel:         getEnvOrDefault("LOG_LEVEL", "info"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	// Validate AI provider configuration
	if cfg.AIProvider != "openrouter" && cfg.AIProvider != "gemini" {
		return nil, fmt.Errorf("AI_PROVIDER must be either 'openrouter' or 'gemini'")
	}

	if cfg.AIProvider == "openrouter" && cfg.OpenRouterAPIKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY is required when AI_PROVIDER is 'openrouter'")
	}

	if cfg.AIProvider == "gemini" && cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required when AI_PROVIDER is 'gemini'")
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
