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
	client := testutils.SetupTestClient(t)

	email := testutils.UniqueEmail("testuserreg")
	password := "testpassword123"

	// Test successful registration.
	resp, err := client.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)

	// Test duplicate registration.
	_, err = client.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err) // Should fail with duplicate email.
}

func TestUserLogin(t *testing.T) {
	client := testutils.SetupTestClient(t)

	email := testutils.UniqueEmail("testuserlogin")
	password := "testpassword123"

	// Register user first.
	_, err := client.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	// Test successful login
	resp, err := client.Login(context.Background(), &clientapi.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)

	// Test login with wrong password
	_, err = client.Login(context.Background(), &clientapi.LoginRequest{
		Email:    email,
		Password: "wrongpassword",
	})
	require.Error(t, err)

	// Test login with non-existent user
	_, err = client.Login(context.Background(), &clientapi.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: password,
	})
	require.Error(t, err)
}

func TestTokenValidation(t *testing.T) {
	client := testutils.SetupTestClient(t)

	email := testutils.UniqueEmail("testtokenval")
	password := "testpassword123"

	// Register and login to get token
	_, err := client.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	loginResp, err := client.Login(context.Background(), &clientapi.LoginRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	token := loginResp.Token

	// Test valid token with data operation
	_, err = client.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: token,
		Data: &clientapi.Data{
			Type:      "test",
			Payload:   []byte("test data"),
			Metadata:  map[string]string{},
			Timestamp: time.Now().Unix(),
		},
	})
	require.NoError(t, err)

	// Test invalid token
	_, err = client.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: "invalid-token",
		Data: &clientapi.Data{
			Type:      "test",
			Payload:   []byte("test data"),
			Metadata:  map[string]string{},
			Timestamp: time.Now().Unix(),
		},
	})
	require.Error(t, err)
}
