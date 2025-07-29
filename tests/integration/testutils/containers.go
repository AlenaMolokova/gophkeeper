// Package testutils provides Testcontainers-based infrastructure for integration testing.
// It manages PostgreSQL containers and generates TLS certificates dynamically
// to ensure test isolation and eliminate hardcoded dependencies.
package testutils

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestEnvironment holds all test infrastructure components.
type TestEnvironment struct {
	PostgresContainer *postgres.PostgresContainer
	CertPath          string
	DSN               string
	Cleanup           func()
}

// SetupTestEnvironment creates an isolated test environment with PostgreSQL and TLS certificates.
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	ctx := context.Background()

	// Create PostgreSQL container
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("gophkeeper_test"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections"),
		),
	)
	require.NoError(t, err, "Failed to start PostgreSQL container")

	// Get database connection string
	dsn, err := postgresContainer.ConnectionString(ctx)
	require.NoError(t, err, "Failed to get database connection string")

	// Generate TLS certificates for testing
	certPath, err := generateTestCertificates(t)
	require.NoError(t, err, "Failed to generate test certificates")

	// Setup cleanup function
	cleanup := func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate PostgreSQL container: %v", err)
		}
		if err := os.RemoveAll(filepath.Dir(certPath)); err != nil {
			t.Logf("Failed to remove test certificates: %v", err)
		}
	}

	return &TestEnvironment{
		PostgresContainer: postgresContainer,
		CertPath:          certPath,
		DSN:               dsn,
		Cleanup:           cleanup,
	}
}

// generateTestCertificates creates temporary TLS certificates for testing.
func generateTestCertificates(t *testing.T) (string, error) {
	// Create temporary directory for certificates
	tempDir := t.TempDir()
	certPath := filepath.Join(tempDir, "server.crt")
	keyPath := filepath.Join(tempDir, "server.key")

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"GophKeeper Test"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{"localhost"},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create certificate: %w", err)
	}

	// Write certificate to file
	certOut, err := os.Create(certPath)
	if err != nil {
		return "", fmt.Errorf("failed to create certificate file: %w", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return "", fmt.Errorf("failed to encode certificate: %w", err)
	}

	// Write private key to file
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return "", fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyOut.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(keyOut, privateKeyPEM); err != nil {
		return "", fmt.Errorf("failed to encode private key: %w", err)
	}

	return certPath, nil
}
