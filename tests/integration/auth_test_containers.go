// Package integration provides Testcontainers-based integration tests for the GophKeeper application.
// These tests use isolated PostgreSQL containers and dynamic TLS certificates
// to ensure complete test isolation and eliminate external dependencies.
package integration

import (
	"context"
	"testing"
	"time"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/AlenaMolokova/gophkeeper/tests/integration/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserRegistrationWithContainers tests user registration using Testcontainers.
func TestUserRegistrationWithContainers(t *testing.T) {
	// Setup isolated test environment
	env := testutils.SetupTestEnvironment(t)
	defer env.Cleanup()

	// Setup test server with containerized database
	testutils.SetupTestServer(t, env)

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	// Setup test clients
	clients := testutils.SetupTestClients(t, env)

	// Test user registration
	ctx := context.Background()
	req := &clientapi.RegisterRequest{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	resp, err := clients.UserClient.Register(ctx, req)
	require.NoError(t, err, "User registration should succeed")
	assert.NotEmpty(t, resp.Token, "Registration should return a token")
}

// TestUserLoginWithContainers tests user login using Testcontainers.
func TestUserLoginWithContainers(t *testing.T) {
	// Setup isolated test environment
	env := testutils.SetupTestEnvironment(t)
	defer env.Cleanup()

	// Setup test server with containerized database
	testutils.SetupTestServer(t, env)

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	// Setup test clients
	clients := testutils.SetupTestClients(t, env)

	// First register a user
	ctx := context.Background()
	registerReq := &clientapi.RegisterRequest{
		Email:    "login@example.com",
		Password: "loginpassword",
	}

	_, err := clients.UserClient.Register(ctx, registerReq)
	require.NoError(t, err, "User registration should succeed")

	// Then test login
	loginReq := &clientapi.LoginRequest{
		Email:    "login@example.com",
		Password: "loginpassword",
	}

	loginResp, err := clients.UserClient.Login(ctx, loginReq)
	require.NoError(t, err, "User login should succeed")
	assert.NotEmpty(t, loginResp.Token, "Login should return a token")
}
