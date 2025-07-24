package integration

import (
	"context"
	"testing"
	"time"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/AlenaMolokova/gophkeeper/tests/integration/testutils"
	"github.com/stretchr/testify/require"
)

func TestDataCRUDOperations(t *testing.T) {
	client := testutils.SetupTestClient(t)

	// Setup: register and login user.
	email := testutils.UniqueEmail("testdatacrud")
	password := "testpassword123"

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

	var dataID string

	// Test AddData.
	t.Run("AddData", func(t *testing.T) {
		data := &clientapi.Data{
			Type:      "login",
			Payload:   []byte("test-password"),
			Metadata:  map[string]string{"url": "https://example.com"},
			Timestamp: time.Now().Unix(),
		}

		resp, err := client.AddData(context.Background(), &clientapi.AddDataRequest{
			Token: token,
			Data:  data,
		})
		require.NoError(t, err)
		require.NotEmpty(t, resp.Id)

		// Store ID for later tests.
		dataID = resp.Id
	})

	// Test GetData.
	t.Run("GetData", func(t *testing.T) {
		require.NotEmpty(t, dataID)

		resp, err := client.GetData(context.Background(), &clientapi.GetDataRequest{
			Token: token,
			Id:    dataID,
		})
		require.NoError(t, err)
		require.Equal(t, dataID, resp.Data.Id)
		require.Equal(t, "login", resp.Data.Type)
		require.Equal(t, []byte("test-password"), resp.Data.Payload)
		require.Equal(t, "https://example.com", resp.Data.Metadata["url"])
	})

	// Test EditData.
	t.Run("EditData", func(t *testing.T) {
		require.NotEmpty(t, dataID)

		updatedData := &clientapi.Data{
			Id:        dataID,
			Type:      "login",
			Payload:   []byte("updated-password"),
			Metadata:  map[string]string{"url": "https://updated-example.com"},
			Timestamp: time.Now().Unix(),
		}

		resp, err := client.EditData(context.Background(), &clientapi.EditDataRequest{
			Token: token,
			Data:  updatedData,
		})
		require.NoError(t, err)
		require.Equal(t, dataID, resp.Id)

		// Verify the update.
		getResp, err := client.GetData(context.Background(), &clientapi.GetDataRequest{
			Token: token,
			Id:    dataID,
		})
		require.NoError(t, err)
		require.Equal(t, []byte("updated-password"), getResp.Data.Payload)
		require.Equal(t, "https://updated-example.com", getResp.Data.Metadata["url"])
	})

	// Test DeleteData.
	t.Run("DeleteData", func(t *testing.T) {
		require.NotEmpty(t, dataID)

		_, err := client.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
			Token: token,
			Id:    dataID,
		})
		require.NoError(t, err)

		// Verify deletion.
		_, err = client.GetData(context.Background(), &clientapi.GetDataRequest{
			Token: token,
			Id:    dataID,
		})
		require.Error(t, err) // Should fail because data is deleted.
	})
}

func TestDataAccessControl(t *testing.T) {
	client := testutils.SetupTestClient(t)

	// Setup: create two users.
	user1Email := testutils.UniqueEmail("user1")
	user2Email := testutils.UniqueEmail("user2")
	password := "testpassword123"

	// Register and login user1.
	_, err := client.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    user1Email,
		Password: password,
	})
	require.NoError(t, err)

	user1Login, err := client.Login(context.Background(), &clientapi.LoginRequest{
		Email:    user1Email,
		Password: password,
	})
	require.NoError(t, err)

	// Register and login user2.
	_, err = client.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    user2Email,
		Password: password,
	})
	require.NoError(t, err)

	user2Login, err := client.Login(context.Background(), &clientapi.LoginRequest{
		Email:    user2Email,
		Password: password,
	})
	require.NoError(t, err)

	// User1 creates data.
	data := &clientapi.Data{
		Type:      "secret",
		Payload:   []byte("user1-secret"),
		Metadata:  map[string]string{},
		Timestamp: time.Now().Unix(),
	}

	addResp, err := client.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: user1Login.Token,
		Data:  data,
	})
	require.NoError(t, err)
	dataID := addResp.Id

	// User2 tries to access user1's data (should fail).
	_, err = client.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user2Login.Token,
		Id:    dataID,
	})
	require.Error(t, err) // Should fail - access denied.

	// User2 tries to edit user1's data (should fail).
	_, err = client.EditData(context.Background(), &clientapi.EditDataRequest{
		Token: user2Login.Token,
		Data: &clientapi.Data{
			Id:        dataID,
			Type:      "secret",
			Payload:   []byte("hacked"),
			Metadata:  map[string]string{},
			Timestamp: time.Now().Unix(),
		},
	})
	require.Error(t, err) // Should fail - access denied.

	// User2 tries to delete user1's data (should fail).
	_, err = client.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: user2Login.Token,
		Id:    dataID,
	})
	require.Error(t, err) // Should fail - access denied.

	// User1 can still access their own data.
	_, err = client.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user1Login.Token,
		Id:    dataID,
	})
	require.NoError(t, err) // Should succeed.
}

func TestDataTypes(t *testing.T) {
	client := testutils.SetupTestClient(t)

	// Setup: register and login user.
	email := testutils.UniqueEmail("testtypes")
	password := "testpassword123"

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

	testCases := []struct {
		name     string
		dataType string
		payload  []byte
		metadata map[string]string
	}{
		{
			name:     "Login",
			dataType: "login",
			payload:  []byte("username:password"),
			metadata: map[string]string{"url": "https://login.example.com"},
		},
		{
			name:     "Text",
			dataType: "text",
			payload:  []byte("Important note"),
			metadata: map[string]string{"category": "notes"},
		},
		{
			name:     "Binary",
			dataType: "binary",
			payload:  []byte{0x01, 0x02, 0x03, 0x04},
			metadata: map[string]string{"filename": "file.bin"},
		},
		{
			name:     "Card",
			dataType: "card",
			payload:  []byte("1234-5678-9012-3456"),
			metadata: map[string]string{"type": "credit"},
		},
		{
			name:     "OTP",
			dataType: "otp",
			payload:  []byte("123456"),
			metadata: map[string]string{"issuer": "Google"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := &clientapi.Data{
				Type:      tc.dataType,
				Payload:   tc.payload,
				Metadata:  tc.metadata,
				Timestamp: time.Now().Unix(),
			}

			// Add data.
			addResp, err := client.AddData(context.Background(), &clientapi.AddDataRequest{
				Token: token,
				Data:  data,
			})
			require.NoError(t, err)
			require.NotEmpty(t, addResp.Id)

			// Get data.
			getResp, err := client.GetData(context.Background(), &clientapi.GetDataRequest{
				Token: token,
				Id:    addResp.Id,
			})
			require.NoError(t, err)
			require.Equal(t, tc.dataType, getResp.Data.Type)
			require.Equal(t, tc.payload, getResp.Data.Payload)
			require.Equal(t, tc.metadata, getResp.Data.Metadata)

			// Clean up.
			_, err = client.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
				Token: token,
				Id:    addResp.Id,
			})
			require.NoError(t, err)
		})
	}
}
