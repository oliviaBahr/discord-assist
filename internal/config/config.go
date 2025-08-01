package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the bot
type Config struct {
	Discord struct {
		Token string
	}
	Bot struct {
		Prefix       string
		Activity     string
		ActivityType string
	}
	Anthropic struct {
		APIKey string
		Model  string
	}
	Server struct {
		Port string
		Host string
	}
	Logging struct {
		Level  string
		Format string
	}
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist in production
	}

	config := &Config{}

	// Discord configuration
	config.Discord.Token = getEnv("DISCORD_TOKEN", "")
	if config.Discord.Token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
	}

	// Bot configuration
	config.Bot.Prefix = getEnv("BOT_PREFIX", "!")
	config.Bot.Activity = getEnv("BOT_ACTIVITY", "with Discord")
	config.Bot.ActivityType = getEnv("BOT_ACTIVITY_TYPE", "Playing")

	// Anthropic configuration
	config.Anthropic.APIKey = getEnv("ANTHROPIC_API_KEY", "")
	if config.Anthropic.APIKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is required")
	}
	config.Anthropic.Model = getEnv("ANTHROPIC_MODEL", "claude-3-5-sonnet")

	// Server configuration
	config.Server.Port = getEnv("SERVER_PORT", "8080")
	config.Server.Host = getEnv("SERVER_HOST", "localhost")

	// Logging configuration
	config.Logging.Level = getEnv("LOG_LEVEL", "info")
	config.Logging.Format = getEnv("LOG_FORMAT", "json")

	return config, nil
}

// getEnv gets an environment variable with a fallback default
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvInt gets an environment variable as an integer with a fallback default
func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getEnvDuration gets an environment variable as a duration with a fallback default
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}
