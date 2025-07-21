// tests/integration/testutils/testutils.go
package testutils

import (
	"os"
	"path/filepath"
	"testing"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// SetupTestClient initializes a gRPC client for integration tests with TLS.
func SetupTestClient(t *testing.T) clientapi.GophKeeperClient {
	certPath := getCertPath(t)
	t.Logf("Loading TLS certificate from: %s", certPath)
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	require.NoError(t, err, "Failed to load TLS certificate")
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	require.NoError(t, err, "Failed to connect to server")
	t.Cleanup(func() { conn.Close() })
	return clientapi.NewGophKeeperClient(conn)
}

// getCertPath returns the path to the TLS certificate.
func getCertPath(t *testing.T) string {
	// Try environment variable first.
	if certPath := os.Getenv("GOPHKEEPER_CERT"); certPath != "" {
		if _, err := os.Stat(certPath); err == nil {
			return filepath.Clean(certPath)
		}
		t.Logf("Certificate not found at GOPHKEEPER_CERT path: %s", certPath)
	}

	// Fallback to project root (E:\go\gophkeeper\cert\server.crt)
	projectRoot := filepath.Join("E:", "go", "gophkeeper")
	certPath := filepath.Join(projectRoot, "cert", "server.crt")
	if _, err := os.Stat(certPath); err == nil {
		return certPath
	}

	t.Fatalf("Certificate not found at %s", certPath)
	return certPath // Unreachable, but required for compilation.
}
