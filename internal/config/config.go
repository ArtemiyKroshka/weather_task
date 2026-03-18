package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	WeatherAPIKey string

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string

	BaseURL    string
	ServerAddr string
}

// Load reads configuration from the environment (.env file is loaded if present).
// Returns an error if any required variable is missing.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:       getEnv("DB_HOST", "db"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "weather_user"),
		DBPassword:   getEnv("DB_PASSWORD", "weather_pass"),
		DBName:       getEnv("DB_NAME", "weather_db"),
		WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		BaseURL:      getEnv("BASE_URL", "http://localhost:8080"),
		ServerAddr:   getEnv("SERVER_ADDR", ":8080"),
	}

	if cfg.WeatherAPIKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY is required")
	}

	port, err := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}
	cfg.SMTPPort = port

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
