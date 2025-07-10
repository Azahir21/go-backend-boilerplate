package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration from environment variables
type Config struct {
    // Database
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string

    // JWT
    JWTSecret     string
    JWTExpiryHours int

    // Server
    ServerPort string
    ServerEnv  string

    // Default Admin
    DefaultAdminUsername string
    DefaultAdminEmail    string
    DefaultAdminPassword string
}

// LoadConfig loads configuration from .env file
func LoadConfig(log *logrus.Logger) *Config {
    err := godotenv.Load()
    if err != nil {
        log.Warn("Error loading .env file, using environment variables")
    }

    expiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "72"))

    return &Config{
        // Database
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "123456"),
        DBName:     getEnv("DB_NAME", "headcount_checker"),
        DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

        // JWT
        JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
        JWTExpiryHours: expiryHours,

        // Server
        ServerPort: getEnv("SERVER_PORT", "8080"),
        ServerEnv:  getEnv("SERVER_ENV", "development"),

        // Default Admin
        DefaultAdminUsername: getEnv("DEFAULT_ADMIN_USERNAME", "admin"),
        DefaultAdminEmail:    getEnv("DEFAULT_ADMIN_EMAIL", "admin@example.com"),
        DefaultAdminPassword: getEnv("DEFAULT_ADMIN_PASSWORD", "admin123"),
    }
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}