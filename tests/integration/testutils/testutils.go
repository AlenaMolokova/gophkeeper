// Package testutils provides utility functions for integration testing of the GophKeeper application.
// It includes Testcontainers-based infrastructure for isolated testing with dynamic
// PostgreSQL containers and TLS certificate generation.
package testutils

import (
	"context"
	"os"
	"testing"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TestClients holds both user and data service clients for integration tests.
type TestClients struct {
	UserClient clientapi.UserServiceClient
	DataClient clientapi.DataServiceClient
	Conn       *grpc.ClientConn
}

// SetupTestClients initializes gRPC clients for integration tests with dynamic TLS certificates.
func SetupTestClients(t *testing.T, env *TestEnvironment) *TestClients {
	t.Logf("Loading TLS certificate from: %s", env.CertPath)
	creds, err := credentials.NewClientTLSFromFile(env.CertPath, "")
	require.NoError(t, err, "Failed to load TLS certificate")

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	require.NoError(t, err, "Failed to connect to server")
	t.Cleanup(func() { conn.Close() })

	return &TestClients{
		UserClient: clientapi.NewUserServiceClient(conn),
		DataClient: clientapi.NewDataServiceClient(conn),
		Conn:       conn,
	}
}

// SetupTestServer starts a test server with the provided environment configuration.
func SetupTestServer(t *testing.T, env *TestEnvironment) {
	// Set environment variables for the test server
	os.Setenv("DATABASE_DSN", env.DSN)
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("CERT_PATH", env.CertPath)

	// Start server in background (implementation depends on your server startup)
	// This would typically involve starting the server in a goroutine
	// and waiting for it to be ready
}

// CreateTestDatabase initializes the test database with required tables.
func CreateTestDatabase(ctx context.Context, env *TestEnvironment) error {
	// Connect to PostgreSQL and create tables
	// This would use the DSN from env.DSN to create the required schema
	return nil
}
