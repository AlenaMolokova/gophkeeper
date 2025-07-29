// Package testutils provides adaptive testing infrastructure that automatically
// chooses between Testcontainers and fallback based on Docker availability.
package testutils

import (
	"fmt"
	"os"
	"testing"
)

const (
	forceFallbackEnv = "true"
)

// TestEnvironmentInterface defines the interface for test environments.
type TestEnvironmentInterface interface {
	GetCertPath() string
	GetDSN() string
	GetCleanup() func()
}

// SetupTestEnvironmentAdaptive automatically chooses the best testing approach.
// It tries Testcontainers first, then falls back to local testing if Docker is not available.
func SetupTestEnvironmentAdaptive(t *testing.T) TestEnvironmentInterface {
	// Check if we should force fallback mode
	if os.Getenv("FORCE_FALLBACK") == forceFallbackEnv {
		t.Log("Using fallback mode (forced)")
		return SetupTestEnvironmentFallback(t)
	}

	// Try to use Testcontainers
	env, err := tryTestcontainers(t)
	if err != nil {
		t.Logf("Testcontainers not available, using fallback: %v", err)
		return SetupTestEnvironmentFallback(t)
	}

	t.Log("Using Testcontainers for isolated testing")
	return env
}

// SetupTestClientsAdaptive automatically chooses the best client setup approach.
func SetupTestClientsAdaptive(t *testing.T, env TestEnvironmentInterface) *TestClients {
	switch e := env.(type) {
	case *TestEnvironment:
		return SetupTestClients(t, e)
	case *TestEnvironmentFallback:
		return SetupTestClientsFallback(t, e)
	default:
		t.Fatalf("Unknown environment type: %T", env)
		return nil
	}
}

// SetupTestServerAdaptive automatically chooses the best server setup approach.
func SetupTestServerAdaptive(t *testing.T, env TestEnvironmentInterface) {
	switch e := env.(type) {
	case *TestEnvironment:
		SetupTestServer(t, e)
	case *TestEnvironmentFallback:
		SetupTestServerFallback(t, e)
	default:
		t.Fatalf("Unknown environment type: %T", env)
	}
}

// tryTestcontainers attempts to create a Testcontainers environment.
// Returns the environment if successful, or an error if Docker is not available.
func tryTestcontainers(t *testing.T) (*TestEnvironment, error) {
	// This is a simplified check - in practice, you might want to
	// actually try to create a container to verify Docker is working
	if os.Getenv("DOCKER_HOST") == "" && os.Getenv("DOCKER_SOCKET") == "" {
		// Try to detect if Docker is available
		if !isDockerAvailable() {
			return nil, fmt.Errorf("Docker not available")
		}
	}

	// Try to create the environment
	return SetupTestEnvironment(t), nil
}

// isDockerAvailable checks if Docker is available on the system.
func isDockerAvailable() bool {
	// Check if Docker Desktop is running (Windows)
	if os.Getenv("DOCKER_DESKTOP") == "true" {
		return true
	}

	// Check if Docker daemon is accessible
	// This is a simplified check - in practice, you might want to
	// actually try to run a Docker command
	return false
}

// GetCertPath returns the certificate path from the environment.
func (e *TestEnvironment) GetCertPath() string {
	return e.CertPath
}

// GetDSN returns the database connection string from the environment.
func (e *TestEnvironment) GetDSN() string {
	return e.DSN
}

// GetCleanup returns the cleanup function from the environment.
func (e *TestEnvironment) GetCleanup() func() {
	return e.Cleanup
}

// GetCertPath returns the certificate path from the fallback environment.
func (e *TestEnvironmentFallback) GetCertPath() string {
	return e.CertPath
}

// GetDSN returns the database connection string from the fallback environment.
func (e *TestEnvironmentFallback) GetDSN() string {
	return e.DSN
}

// GetCleanup returns the cleanup function from the fallback environment.
func (e *TestEnvironmentFallback) GetCleanup() func() {
	return e.Cleanup
}
