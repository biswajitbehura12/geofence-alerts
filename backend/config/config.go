package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "geofence"),
		},
	}
}

// GetDatabaseURL returns the database connection string
func (c *DatabaseConfig) GetConnectionURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
