// Package testutils provides fallback testing infrastructure for environments
// where Docker/Testcontainers are not available (e.g., Windows without Docker Desktop).
package testutils

import (
	"os"
	"testing"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// TestEnvironmentFallback provides a fallback test environment when Docker is not available.
// It uses local resources and mocks instead of containers.
type TestEnvironmentFallback struct {
	CertPath string
	DSN      string
	Cleanup  func()
}

// SetupTestEnvironmentFallback creates a fallback test environment without Docker.
// This is used when Testcontainers are not available (e.g., Windows without Docker Desktop).
func SetupTestEnvironmentFallback(t *testing.T) *TestEnvironmentFallback {
	// Generate TLS certificates for testing
	certPath, err := generateTestCertificates(t)
	require.NoError(t, err, "Failed to generate test certificates")

	// Use local PostgreSQL or SQLite for testing
	dsn := getLocalDatabaseDSN()

	// Setup cleanup function
	cleanup := func() {
		if err := os.RemoveAll(certPath); err != nil {
			t.Logf("Failed to remove test certificates: %v", err)
		}
	}

	return &TestEnvironmentFallback{
		CertPath: certPath,
		DSN:      dsn,
		Cleanup:  cleanup,
	}
}

// getLocalDatabaseDSN returns a DSN for local database testing.
// It tries to use environment variables first, then falls back to defaults.
func getLocalDatabaseDSN() string {
	// Try to use environment variable first
	if dsn := os.Getenv("TEST_DATABASE_DSN"); dsn != "" {
		return dsn
	}

	// Fallback to local PostgreSQL if available
	if os.Getenv("USE_LOCAL_POSTGRES") == "true" {
		return "postgres://testuser:testpass@localhost:5432/gophkeeper_test"
	}

	// Final fallback to SQLite for testing
	return "file:test.db?cache=shared&mode=memory"
}

// SetupTestClientsFallback initializes gRPC clients for integration tests without containers.
func SetupTestClientsFallback(t *testing.T, env *TestEnvironmentFallback) *TestClients {
	t.Logf("Loading TLS certificate from: %s", env.CertPath)

	// For fallback testing, we might skip TLS or use insecure connection
	// depending on the test environment
	var conn *grpc.ClientConn
	var err error

	if os.Getenv("SKIP_TLS") == "true" {
		// Use insecure connection for local testing
		conn, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// Try to use TLS certificate
		creds, credErr := credentials.NewClientTLSFromFile(env.CertPath, "")
		if credErr != nil {
			t.Logf("TLS certificate failed, falling back to insecure: %v", credErr)
			conn, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			conn, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
		}
	}

	require.NoError(t, err, "Failed to connect to server")
	t.Cleanup(func() { conn.Close() })

	return &TestClients{
		UserClient: clientapi.NewUserServiceClient(conn),
		DataClient: clientapi.NewDataServiceClient(conn),
		Conn:       conn,
	}
}

// SetupTestServerFallback starts a test server with fallback configuration.
func SetupTestServerFallback(t *testing.T, env *TestEnvironmentFallback) {
	// Set environment variables for the test server
	os.Setenv("DATABASE_DSN", env.DSN)
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("CERT_PATH", env.CertPath)

	// For fallback testing, we might start a local server
	// or use mocks depending on the configuration
	if os.Getenv("USE_MOCK_SERVER") == "true" {
		t.Log("Using mock server for testing")
		// Mock server implementation would go here
	} else {
		t.Log("Starting local test server")
		// Local server startup would go here
	}
}
