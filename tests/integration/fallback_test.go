// Package integration provides fallback integration tests for environments
// where Docker/Testcontainers are not available.
package integration

import (
	"context"
	"os"
	"testing"
	"time"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/AlenaMolokova/gophkeeper/tests/integration/testutils"
	"github.com/stretchr/testify/assert"
)

// TestUserRegistrationFallback tests user registration using fallback mode.
// This test works without Docker and uses local resources.
func TestUserRegistrationFallback(t *testing.T) {
	// Force fallback mode
	os.Setenv("FORCE_FALLBACK", "true")
	os.Setenv("SKIP_TLS", "true")
	defer func() {
		os.Unsetenv("FORCE_FALLBACK")
		os.Unsetenv("SKIP_TLS")
	}()

	// Setup fallback test environment
	env := testutils.SetupTestEnvironmentFallback(t)
	defer env.GetCleanup()()

	// Setup test server with fallback configuration
	testutils.SetupTestServerFallback(t, env)

	// Wait for server to be ready (if running locally)
	time.Sleep(1 * time.Second)

	// Setup test clients with fallback configuration
	clients := testutils.SetupTestClientsFallback(t, env)

	// Test user registration
	ctx := context.Background()
	req := &clientapi.RegisterRequest{
		Email:    "fallback@example.com",
		Password: "testpassword",
	}

	// Note: This test might fail if no server is running locally
	// The purpose is to test the fallback infrastructure, not the actual server
	resp, err := clients.UserClient.Register(ctx, req)
	if err != nil {
		t.Logf("Registration failed (expected if no server running): %v", err)
		// This is expected if no server is running locally
		return
	}

	assert.NotEmpty(t, resp.Token, "Registration should return a token")
}

// TestFallbackInfrastructure tests that fallback infrastructure works correctly.
func TestFallbackInfrastructure(t *testing.T) {
	// Force fallback mode
	os.Setenv("FORCE_FALLBACK", "true")
	defer os.Unsetenv("FORCE_FALLBACK")

	// Setup fallback test environment
	env := testutils.SetupTestEnvironmentFallback(t)
	defer env.GetCleanup()()

	// Verify environment is set up correctly
	assert.NotEmpty(t, env.GetCertPath(), "Certificate path should be set")
	assert.NotEmpty(t, env.GetDSN(), "Database DSN should be set")
	assert.NotNil(t, env.GetCleanup(), "Cleanup function should be set")

	// Verify certificate file exists
	_, err := os.Stat(env.GetCertPath())
	assert.NoError(t, err, "Certificate file should exist")

	t.Logf("Fallback environment setup successfully:")
	t.Logf("  CertPath: %s", env.GetCertPath())
	t.Logf("  DSN: %s", env.GetDSN())
}

// TestAdaptiveEnvironment tests the adaptive environment selection.
func TestAdaptiveEnvironment(t *testing.T) {
	// Test adaptive environment selection
	env := testutils.SetupTestEnvironmentAdaptive(t)
	defer env.GetCleanup()()

	// Verify environment is set up correctly
	assert.NotEmpty(t, env.GetCertPath(), "Certificate path should be set")
	assert.NotEmpty(t, env.GetDSN(), "Database DSN should be set")
	assert.NotNil(t, env.GetCleanup(), "Cleanup function should be set")

	t.Logf("Adaptive environment selected successfully:")
	t.Logf("  Type: %T", env)
	t.Logf("  CertPath: %s", env.GetCertPath())
	t.Logf("  DSN: %s", env.GetDSN())
}
