// Package config provides centralized configuration management for the GophKeeper application.
// It handles both server and client configurations with environment variable support,
// default values, and validation.
//
// The package follows Go best practices by providing type-safe configuration
// structures and helper functions for parsing environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// ServerConfig holds server configuration.
// It contains all necessary settings for running the GophKeeper server,
// including database, JWT, TLS, and server network settings.
type ServerConfig struct {
	Database DatabaseConfig
	JWT      JWTConfig
	TLS      TLSConfig
	Server   ServerSettings
}

// DatabaseConfig holds database configuration.
// It includes connection string, pool settings, and health check parameters
// for optimal database performance and reliability.
type DatabaseConfig struct {
	DSN               string        // Database connection string
	MaxConns          int32         // Maximum number of connections in the pool
	MinConns          int32         // Minimum number of connections to maintain
	MaxConnLifetime   time.Duration // Maximum time a connection can be reused
	MaxConnIdleTime   time.Duration // Maximum time a connection can remain idle
	HealthCheckPeriod time.Duration // Interval between health checks
}

// JWTConfig holds JWT configuration.
// It contains the secret key for token signing and token lifetime settings.
type JWTConfig struct {
	Secret string        // Secret key for JWT token signing
	TTL    time.Duration // Token time-to-live duration
}

// TLSConfig holds TLS configuration.
// It specifies the paths to certificate and private key files for secure communication.
type TLSConfig struct {
	CertFile string // Path to the TLS certificate file
	KeyFile  string // Path to the TLS private key file
}

// ServerSettings holds server settings.
// It defines the network interface and port for the server to listen on.
type ServerSettings struct {
	Host string // Server host address
	Port int    // Server port number
}

// ClientConfig holds client configuration.
// It contains settings needed for client applications to connect to the server.
type ClientConfig struct {
	Token    string         // Authentication token for client requests
	CertPath string         // Path to the server TLS certificate
	Server   ServerSettings // Server connection settings
}

// getEnvAsInt32OrDefault retrieves an environment variable as an int32 or returns a default.
// It attempts to parse the environment variable as an integer and returns
// the default value if parsing fails or the variable is not set.
// The function safely handles integer overflow by clamping values to int32 range.
func getEnvAsInt32OrDefault(key string, defaultValue int32) int32 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			// Safely convert to int32, clamping to valid range
			if intValue > 0 && intValue <= 2147483647 { // max int32
				//nolint:gosec // We've already validated the range
				return int32(intValue)
			}
		}
	}
	return defaultValue
}

// NewServerConfig creates a new ServerConfig from environment variables.
// It reads all necessary configuration values from environment variables
// with sensible defaults and performs validation to ensure required values are present.
//
// Required environment variables:
//   - DATABASE_DSN: PostgreSQL connection string
//   - JWT_SECRET: Secret key for JWT token signing
//
// Optional environment variables with defaults:
//   - DB_MAX_CONNS, DB_MIN_CONNS, DB_MAX_CONN_LIFETIME, etc.
//   - JWT_TTL, TLS_CERT_FILE, TLS_KEY_FILE, SERVER_HOST, SERVER_PORT
func NewServerConfig() (*ServerConfig, error) {
	config := &ServerConfig{
		Database: DatabaseConfig{
			DSN:               getEnvOrDefault("DATABASE_DSN", ""),
			MaxConns:          getEnvAsInt32OrDefault("DB_MAX_CONNS", 20),
			MinConns:          getEnvAsInt32OrDefault("DB_MIN_CONNS", 5),
			MaxConnLifetime:   getEnvAsDurationOrDefault("DB_MAX_CONN_LIFETIME", 1*time.Hour),
			MaxConnIdleTime:   getEnvAsDurationOrDefault("DB_MAX_CONN_IDLE_TIME", 30*time.Minute),
			HealthCheckPeriod: getEnvAsDurationOrDefault("DB_HEALTH_CHECK_PERIOD", 1*time.Minute),
		},
		JWT: JWTConfig{
			Secret: getEnvOrDefault("JWT_SECRET", ""),
			TTL:    getEnvAsDurationOrDefault("JWT_TTL", 24*time.Hour),
		},
		TLS: TLSConfig{
			CertFile: getEnvOrDefault("TLS_CERT_FILE", "cert/server.crt"),
			KeyFile:  getEnvOrDefault("TLS_KEY_FILE", "cert/server.key"),
		},
		Server: ServerSettings{
			Host: getEnvOrDefault("SERVER_HOST", "localhost"),
			Port: getEnvAsIntOrDefault("SERVER_PORT", 50051),
		},
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// NewClientConfig creates a new ClientConfig from environment variables.
// It provides client-specific configuration with defaults suitable for development
// and testing environments.
//
// Environment variables:
//   - GOPHKEEPER_TOKEN: Authentication token (default: "test-token")
//   - GOPHKEEPER_CERT: Path to server certificate (default: "cert/server.crt")
//   - SERVER_HOST, SERVER_PORT: Server connection settings
func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		Token:    getEnvOrDefault("GOPHKEEPER_TOKEN", "test-token"),
		CertPath: getEnvOrDefault("GOPHKEEPER_CERT", "cert/server.crt"),
		Server: ServerSettings{
			Host: getEnvOrDefault("SERVER_HOST", "localhost"),
			Port: getEnvAsIntOrDefault("SERVER_PORT", 50051),
		},
	}
}

// validate validates the server configuration.
// It ensures that all required configuration values are present
// and returns an error if any required values are missing.
func (c *ServerConfig) validate() error {
	if c.Database.DSN == "" {
		return fmt.Errorf("DATABASE_DSN environment variable is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	return nil
}

// GetDSN returns the database connection string.
// This method provides a convenient way to access the database DSN
// from the configuration structure.
func (c *ServerConfig) GetDSN() string {
	return c.Database.DSN
}

// GetJWTSecret returns the JWT secret as bytes.
// This method converts the string secret to bytes for use with JWT libraries.
func (c *ServerConfig) GetJWTSecret() []byte {
	return []byte(c.JWT.Secret)
}

// GetServerAddress returns the server address as string.
// It formats the host and port into a standard address string
// suitable for network connections.
func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetClientServerAddress returns the server address for client connections.
// It formats the host and port into a standard address string
// suitable for client network connections.
func (c *ClientConfig) GetClientServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// Helper functions for environment variable parsing

// getEnvOrDefault retrieves an environment variable value or returns a default.
// It checks if the environment variable exists and is not empty,
// returning the default value if the variable is not set or empty.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault retrieves an environment variable as an integer or returns a default.
// It attempts to parse the environment variable as an integer and returns
// the default value if parsing fails or the variable is not set.
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDurationOrDefault retrieves an environment variable as a duration or returns a default.
// It attempts to parse the environment variable as a time.Duration and returns
// the default value if parsing fails or the variable is not set.
func getEnvAsDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
