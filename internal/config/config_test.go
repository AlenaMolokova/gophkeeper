package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServerConfig(t *testing.T) {
	// Set required environment variables
	os.Setenv("DATABASE_DSN", "postgres://test:test@localhost:5432/test")
	os.Setenv("JWT_SECRET", "test-secret")
	defer func() {
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("JWT_SECRET")
	}()

	config, err := NewServerConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	// Test default values
	assert.Equal(t, "postgres://test:test@localhost:5432/test", config.Database.DSN)
	assert.Equal(t, "test-secret", config.JWT.Secret)
	assert.Equal(t, int32(20), config.Database.MaxConns)
	assert.Equal(t, int32(5), config.Database.MinConns)
	assert.Equal(t, 1*time.Hour, config.Database.MaxConnLifetime)
	assert.Equal(t, 30*time.Minute, config.Database.MaxConnIdleTime)
	assert.Equal(t, 1*time.Minute, config.Database.HealthCheckPeriod)
	assert.Equal(t, 24*time.Hour, config.JWT.TTL)
	assert.Equal(t, "cert/server.crt", config.TLS.CertFile)
	assert.Equal(t, "cert/server.key", config.TLS.KeyFile)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, 50051, config.Server.Port)
}

func TestNewServerConfigWithCustomValues(t *testing.T) {
	// Set all environment variables
	os.Setenv("DATABASE_DSN", "postgres://custom:custom@localhost:5432/custom")
	os.Setenv("JWT_SECRET", "custom-secret")
	os.Setenv("DB_MAX_CONNS", "50")
	os.Setenv("DB_MIN_CONNS", "10")
	os.Setenv("DB_MAX_CONN_LIFETIME", "2h")
	os.Setenv("DB_MAX_CONN_IDLE_TIME", "1h")
	os.Setenv("DB_HEALTH_CHECK_PERIOD", "30s")
	os.Setenv("JWT_TTL", "12h")
	os.Setenv("TLS_CERT_FILE", "custom.crt")
	os.Setenv("TLS_KEY_FILE", "custom.key")
	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("SERVER_PORT", "8080")

	defer func() {
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DB_MAX_CONNS")
		os.Unsetenv("DB_MIN_CONNS")
		os.Unsetenv("DB_MAX_CONN_LIFETIME")
		os.Unsetenv("DB_MAX_CONN_IDLE_TIME")
		os.Unsetenv("DB_HEALTH_CHECK_PERIOD")
		os.Unsetenv("JWT_TTL")
		os.Unsetenv("TLS_CERT_FILE")
		os.Unsetenv("TLS_KEY_FILE")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
	}()

	config, err := NewServerConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	// Test custom values
	assert.Equal(t, "postgres://custom:custom@localhost:5432/custom", config.Database.DSN)
	assert.Equal(t, "custom-secret", config.JWT.Secret)
	assert.Equal(t, int32(50), config.Database.MaxConns)
	assert.Equal(t, int32(10), config.Database.MinConns)
	assert.Equal(t, 2*time.Hour, config.Database.MaxConnLifetime)
	assert.Equal(t, 1*time.Hour, config.Database.MaxConnIdleTime)
	assert.Equal(t, 30*time.Second, config.Database.HealthCheckPeriod)
	assert.Equal(t, 12*time.Hour, config.JWT.TTL)
	assert.Equal(t, "custom.crt", config.TLS.CertFile)
	assert.Equal(t, "custom.key", config.TLS.KeyFile)
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, 8080, config.Server.Port)
}

func TestNewServerConfigValidation(t *testing.T) {
	// Test missing DATABASE_DSN
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("JWT_SECRET")

	config, err := NewServerConfig()
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "DATABASE_DSN environment variable is required")

	// Test missing JWT_SECRET
	os.Setenv("DATABASE_DSN", "postgres://test:test@localhost:5432/test")
	os.Unsetenv("JWT_SECRET")

	config, err = NewServerConfig()
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "JWT_SECRET environment variable is required")
}

func TestNewClientConfig(t *testing.T) {
	// Test default values
	config := NewClientConfig()
	require.NotNil(t, config)

	assert.Equal(t, "test-token", config.Token)
	assert.Equal(t, "cert/server.crt", config.CertPath)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, 50051, config.Server.Port)

	// Test custom values
	os.Setenv("GOPHKEEPER_TOKEN", "custom-token")
	os.Setenv("GOPHKEEPER_CERT", "custom.crt")
	os.Setenv("SERVER_HOST", "custom-host")
	os.Setenv("SERVER_PORT", "9090")

	defer func() {
		os.Unsetenv("GOPHKEEPER_TOKEN")
		os.Unsetenv("GOPHKEEPER_CERT")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
	}()

	config = NewClientConfig()
	require.NotNil(t, config)

	assert.Equal(t, "custom-token", config.Token)
	assert.Equal(t, "custom.crt", config.CertPath)
	assert.Equal(t, "custom-host", config.Server.Host)
	assert.Equal(t, 9090, config.Server.Port)
}

func TestServerConfigMethods(t *testing.T) {
	os.Setenv("DATABASE_DSN", "postgres://test:test@localhost:5432/test")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("SERVER_HOST", "test-host")
	os.Setenv("SERVER_PORT", "1234")

	defer func() {
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
	}()

	config, err := NewServerConfig()
	require.NoError(t, err)

	// Test GetDSN
	assert.Equal(t, "postgres://test:test@localhost:5432/test", config.GetDSN())

	// Test GetJWTSecret
	assert.Equal(t, []byte("test-secret"), config.GetJWTSecret())

	// Test GetServerAddress
	assert.Equal(t, "test-host:1234", config.GetServerAddress())
}

func TestClientConfigMethods(t *testing.T) {
	os.Setenv("SERVER_HOST", "client-host")
	os.Setenv("SERVER_PORT", "5678")

	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
	}()

	config := NewClientConfig()

	// Test GetClientServerAddress
	assert.Equal(t, "client-host:5678", config.GetClientServerAddress())
}

func TestHelperFunctions(t *testing.T) {
	// Test getEnvOrDefault
	os.Setenv("TEST_KEY", "test-value")
	defer os.Unsetenv("TEST_KEY")

	assert.Equal(t, "test-value", getEnvOrDefault("TEST_KEY", "default"))
	assert.Equal(t, "default", getEnvOrDefault("NONEXISTENT_KEY", "default"))

	// Test getEnvAsIntOrDefault
	os.Setenv("INT_KEY", "42")
	defer os.Unsetenv("INT_KEY")

	assert.Equal(t, 42, getEnvAsIntOrDefault("INT_KEY", 0))
	assert.Equal(t, 0, getEnvAsIntOrDefault("NONEXISTENT_INT_KEY", 0))
	assert.Equal(t, 0, getEnvAsIntOrDefault("INVALID_INT_KEY", 0))

	// Test getEnvAsDurationOrDefault
	os.Setenv("DURATION_KEY", "1h30m")
	defer os.Unsetenv("DURATION_KEY")

	expectedDuration := 1*time.Hour + 30*time.Minute
	assert.Equal(t, expectedDuration, getEnvAsDurationOrDefault("DURATION_KEY", time.Hour))
	assert.Equal(t, time.Hour, getEnvAsDurationOrDefault("NONEXISTENT_DURATION_KEY", time.Hour))
	assert.Equal(t, time.Hour, getEnvAsDurationOrDefault("INVALID_DURATION_KEY", time.Hour))
}
