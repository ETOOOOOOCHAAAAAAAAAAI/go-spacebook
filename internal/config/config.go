package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	API      APIConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type ServerConfig struct {
	Port string
	Mode string
}

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type APIConfig struct {
	Version string
	Prefix  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "spacebook_db"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET_KEY", "your-secret-key"),
			AccessTokenTTL:  parseDuration(getEnv("JWT_ACCESS_TOKEN_TTL", "15m"), 15*time.Minute),
			RefreshTokenTTL: parseDuration(getEnv("JWT_REFRESH_TOKEN_TTL", "168h"), 168*time.Hour),
		},
		API: APIConfig{
			Version: getEnv("API_VERSION", "v1"),
			Prefix:  getEnv("API_PREFIX", "/api"),
		},
	}

	return config, nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string, defaultDuration time.Duration) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return defaultDuration
	}
	return duration
}
