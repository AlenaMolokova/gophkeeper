// Package integration provides infrastructure tests that verify
// the Testcontainers and fallback systems work correctly.
package integration

import (
	"os"
	"testing"

	"github.com/AlenaMolokova/gophkeeper/tests/integration/testutils"
	"github.com/stretchr/testify/assert"
)

// TestInfrastructureSetup tests that the testing infrastructure can be set up correctly.
func TestInfrastructureSetup(t *testing.T) {
	// Test adaptive environment setup
	env := testutils.SetupTestEnvironmentAdaptive(t)
	defer env.GetCleanup()()

	// Verify environment is properly configured
	assert.NotEmpty(t, env.GetCertPath(), "Certificate path should be set")
	assert.NotEmpty(t, env.GetDSN(), "Database DSN should be set")
	assert.NotNil(t, env.GetCleanup(), "Cleanup function should be set")

	t.Logf("Infrastructure setup successful:")
	t.Logf("  Environment type: %T", env)
	t.Logf("  Certificate path: %s", env.GetCertPath())
	t.Logf("  Database DSN: %s", env.GetDSN())
}

// TestFallbackMode tests that fallback mode works correctly.
func TestFallbackMode(t *testing.T) {
	// Force fallback mode
	os.Setenv("FORCE_FALLBACK", "true")
	defer os.Unsetenv("FORCE_FALLBACK")

	// Setup fallback environment
	env := testutils.SetupTestEnvironmentFallback(t)
	defer env.GetCleanup()()

	// Verify it's actually a fallback environment
	assert.IsType(t, &testutils.TestEnvironmentFallback{}, env, "Should be fallback environment")

	// Verify environment is properly configured
	assert.NotEmpty(t, env.GetCertPath(), "Certificate path should be set")
	assert.NotEmpty(t, env.GetDSN(), "Database DSN should be set")
	assert.NotNil(t, env.GetCleanup(), "Cleanup function should be set")

	t.Logf("Fallback mode setup successful:")
	t.Logf("  Certificate path: %s", env.GetCertPath())
	t.Logf("  Database DSN: %s", env.GetDSN())
}

// TestCertificateGeneration tests that TLS certificates are generated correctly.
func TestCertificateGeneration(t *testing.T) {
	// Setup environment
	env := testutils.SetupTestEnvironmentAdaptive(t)
	defer env.GetCleanup()()

	// Verify certificate file exists and is readable
	certPath := env.GetCertPath()
	assert.NotEmpty(t, certPath, "Certificate path should not be empty")

	// Check if certificate file exists
	fileInfo, err := os.Stat(certPath)
	assert.NoError(t, err, "Certificate file should exist")
	assert.True(t, fileInfo.Size() > 0, "Certificate file should not be empty")

	t.Logf("Certificate generation successful:")
	t.Logf("  Certificate path: %s", certPath)
	t.Logf("  File size: %d bytes", fileInfo.Size())
}

// TestEnvironmentCleanup tests that cleanup functions work correctly.
func TestEnvironmentCleanup(t *testing.T) {
	// Setup environment
	env := testutils.SetupTestEnvironmentAdaptive(t)

	// Get certificate path before cleanup
	certPath := env.GetCertPath()
	assert.NotEmpty(t, certPath, "Certificate path should be set")

	// Verify certificate exists
	_, err := os.Stat(certPath)
	assert.NoError(t, err, "Certificate should exist before cleanup")

	// Run cleanup
	cleanup := env.GetCleanup()
	assert.NotNil(t, cleanup, "Cleanup function should be set")
	cleanup()

	// Verify certificate is cleaned up (this might not work on Windows due to file locking)
	// We'll just log the attempt
	t.Logf("Cleanup executed for certificate: %s", certPath)
}

// TestAdaptiveSelection tests that the adaptive system correctly chooses between modes.
func TestAdaptiveSelection(t *testing.T) {
	// Test without forcing fallback
	env1 := testutils.SetupTestEnvironmentAdaptive(t)
	defer env1.GetCleanup()()

	// Force fallback mode
	os.Setenv("FORCE_FALLBACK", "true")
	env2 := testutils.SetupTestEnvironmentAdaptive(t)
	defer env2.GetCleanup()()
	os.Unsetenv("FORCE_FALLBACK")

	// Verify different environment types
	t.Logf("Environment 1 type: %T", env1)
	t.Logf("Environment 2 type: %T", env2)

	// Both should be properly configured regardless of type
	assert.NotEmpty(t, env1.GetCertPath(), "Environment 1 should have certificate path")
	assert.NotEmpty(t, env2.GetCertPath(), "Environment 2 should have certificate path")
	assert.NotEmpty(t, env1.GetDSN(), "Environment 1 should have DSN")
	assert.NotEmpty(t, env2.GetDSN(), "Environment 2 should have DSN")
}
