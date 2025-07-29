package integration

import (
	"context"
	"testing"
	"time"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	testutils "github.com/AlenaMolokova/gophkeeper/tests/integration/testutils"
	"github.com/stretchr/testify/require"
)

func TestUserRegistration(t *testing.T) {
	// Setup adaptive test environment (automatically chooses Testcontainers or fallback)
	env := testutils.SetupTestEnvironmentAdaptive(t)
	defer env.GetCleanup()()

	// Setup test server with adaptive configuration
	testutils.SetupTestServerAdaptive(t, env)

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	// Setup test clients with adaptive configuration
	clients := testutils.SetupTestClientsAdaptive(t, env)

	email := testutils.UniqueEmail("testuserreg")
	password := "testpassword123"

	// Test successful registration.
	resp, err := clients.UserClient.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)

	// Test duplicate registration.
	_, err = clients.UserClient.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err) // Should fail with duplicate email.
}

func TestUserLogin(t *testing.T) {
	// Setup adaptive test environment
	env := testutils.SetupTestEnvironmentAdaptive(t)
	defer env.GetCleanup()()

	// Setup test server with adaptive configuration
	testutils.SetupTestServerAdaptive(t, env)

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	// Setup test clients with adaptive configuration
	clients := testutils.SetupTestClientsAdaptive(t, env)

	email := testutils.UniqueEmail("testuserlogin")
	password := "testpassword123"

	// Register user first.
	_, err := clients.UserClient.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	// Test successful login
	resp, err := clients.UserClient.Login(context.Background(), &clientapi.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)

	// Test login with wrong password
	_, err = clients.UserClient.Login(context.Background(), &clientapi.LoginRequest{
		Email:    email,
		Password: "wrongpassword",
	})
	require.Error(t, err)

	// Test login with non-existent user
	_, err = clients.UserClient.Login(context.Background(), &clientapi.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: password,
	})
	require.Error(t, err)
}

func TestTokenValidation(t *testing.T) {
	// Setup adaptive test environment
	env := testutils.SetupTestEnvironmentAdaptive(t)
	defer env.GetCleanup()()

	// Setup test server with adaptive configuration
	testutils.SetupTestServerAdaptive(t, env)

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	// Setup test clients with adaptive configuration
	clients := testutils.SetupTestClientsAdaptive(t, env)

	email := testutils.UniqueEmail("testtokenval")
	password := "testpassword123"

	// Register and login to get token
	_, err := clients.UserClient.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	loginResp, err := clients.UserClient.Login(context.Background(), &clientapi.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	token := loginResp.Token

	// Test valid token with data operation
	_, err = clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: token,
		Data: &clientapi.Data{
			Type:      clientapi.DataType_DATA_TYPE_TEXT,
			Payload:   []byte("test data"),
			Timestamp: time.Now().Unix(),
		},
	})
	require.NoError(t, err)

	// Test invalid token
	_, err = clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: "invalid-token",
		Data: &clientapi.Data{
			Type:      clientapi.DataType_DATA_TYPE_TEXT,
			Payload:   []byte("test data"),
			Timestamp: time.Now().Unix(),
		},
	})
	require.Error(t, err)
}
