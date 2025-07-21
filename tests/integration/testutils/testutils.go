package testutils

import (
	"os"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// SetupTestClient creates a gRPC client for testing.
func SetupTestClient(t require.TestingT) clientapi.GophKeeperClient {
	certPath := getCertPath()
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	require.NoError(t, err)
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	require.NoError(t, err)
	return clientapi.NewGophKeeperClient(conn)
}

// getCertPath returns the path to the TLS certificate.
func getCertPath() string {
	// Try environment variable first.
	if certPath := os.Getenv("GOPHKEEPER_CERT"); certPath != "" {
		return certPath
	}

	// Try relative to current directory.
	if _, err := os.Stat("cert/server.crt"); err == nil {
		return "cert/server.crt"
	}

	// Try relative to project root (when running from tests directory).
	if _, err := os.Stat("../cert/server.crt"); err == nil {
		return "../cert/server.crt"
	}

	// Try relative to project root (when running from tests/integration directory).
	if _, err := os.Stat("../../cert/server.crt"); err == nil {
		return "../../cert/server.crt"
	}

	// Fallback to default path.
	return "cert/server.crt"
}
