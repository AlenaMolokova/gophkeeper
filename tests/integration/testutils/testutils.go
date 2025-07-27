// Package testutils provides utility functions for integration testing of the GophKeeper application.
// It includes helper functions for setting up test clients and managing test resources
// with proper TLS configuration and cleanup for integration test scenarios.
package testutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AlenaMolokova/gophkeeper/internal/config"
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

// SetupTestClients initializes gRPC clients for integration tests with TLS.
func SetupTestClients(t *testing.T) *TestClients {
	certPath := getCertPath(t)
	t.Logf("Loading TLS certificate from: %s", certPath)
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
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

// getCertPath returns the path to the TLS certificate.
func getCertPath(t *testing.T) string {
	// Use config to get certificate path
	clientConfig := config.NewClientConfig()
	certPath := clientConfig.CertPath

	// Check if certificate exists at config path
	if _, err := os.Stat(certPath); err == nil {
		return filepath.Clean(certPath)
	}
	t.Logf("Certificate not found at config path: %s", certPath)

	// Fallback to project root (E:\go\gophkeeper\cert\server.crt)
	projectRoot := filepath.Join("E:", "go", "gophkeeper")
	fallbackCertPath := filepath.Join(projectRoot, "cert", "server.crt")
	if _, err := os.Stat(fallbackCertPath); err == nil {
		return fallbackCertPath
	}

	t.Fatalf("Certificate not found at %s or %s", certPath, fallbackCertPath)
	return certPath // Unreachable, but required for compilation.
}
