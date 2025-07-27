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
	clients := testutils.SetupTestClients(t)

	// Setup: register and login user.
	email := testutils.UniqueEmail("testdatacrud")
	password := "testpassword123"

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

	var dataID string

	// Test AddData.
	t.Run("AddData", func(t *testing.T) {
		data := &clientapi.Data{
			Type:      clientapi.DataType_DATA_TYPE_LOGIN,
			Payload:   []byte("test-password"),
			Metadata:  &clientapi.Data_LoginData{LoginData: &clientapi.LoginData{Url: "https://example.com"}},
			Timestamp: time.Now().Unix(),
		}

		resp, err := clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
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

		resp, err := clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
			Token: token,
			Id:    dataID,
		})
		require.NoError(t, err)
		require.Equal(t, dataID, resp.Data.Id)
		require.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, resp.Data.Type)
		require.Equal(t, []byte("test-password"), resp.Data.Payload)

		// Check metadata
		if loginData := resp.Data.GetLoginData(); loginData != nil {
			require.Equal(t, "https://example.com", loginData.Url)
		}
	})

	// Test EditData.
	t.Run("EditData", func(t *testing.T) {
		require.NotEmpty(t, dataID)

		updatedData := &clientapi.Data{
			Id:        dataID,
			Type:      clientapi.DataType_DATA_TYPE_LOGIN,
			Payload:   []byte("updated-password"),
			Metadata:  &clientapi.Data_LoginData{LoginData: &clientapi.LoginData{Url: "https://updated-example.com"}},
			Timestamp: time.Now().Unix(),
		}

		resp, err := clients.DataClient.EditData(context.Background(), &clientapi.EditDataRequest{
			Token: token,
			Data:  updatedData,
		})
		require.NoError(t, err)
		require.Equal(t, dataID, resp.Id)

		// Verify the update.
		getResp, err := clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
			Token: token,
			Id:    dataID,
		})
		require.NoError(t, err)
		require.Equal(t, []byte("updated-password"), getResp.Data.Payload)

		// Check updated metadata
		if loginData := getResp.Data.GetLoginData(); loginData != nil {
			require.Equal(t, "https://updated-example.com", loginData.Url)
		}
	})

	// Test DeleteData.
	t.Run("DeleteData", func(t *testing.T) {
		require.NotEmpty(t, dataID)

		_, err := clients.DataClient.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
			Token: token,
			Id:    dataID,
		})
		require.NoError(t, err)

		// Verify deletion.
		_, err = clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
			Token: token,
			Id:    dataID,
		})
		require.Error(t, err) // Should fail because data was deleted.
	})
}

func TestDataAccessControl(t *testing.T) {
	clients := testutils.SetupTestClients(t)

	// Setup: create two users.
	user1Email := testutils.UniqueEmail("user1")
	user2Email := testutils.UniqueEmail("user2")
	password := "testpassword123"

	// Register and login user1.
	_, err := clients.UserClient.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    user1Email,
		Password: password,
	})
	require.NoError(t, err)

	user1Login, err := clients.UserClient.Login(context.Background(), &clientapi.LoginRequest{
		Email:    user1Email,
		Password: password,
	})
	require.NoError(t, err)

	// Register and login user2.
	_, err = clients.UserClient.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    user2Email,
		Password: password,
	})
	require.NoError(t, err)

	user2Login, err := clients.UserClient.Login(context.Background(), &clientapi.LoginRequest{
		Email:    user2Email,
		Password: password,
	})
	require.NoError(t, err)

	// User1 creates data.
	data := &clientapi.Data{
		Type:      clientapi.DataType_DATA_TYPE_TEXT,
		Payload:   []byte("user1-secret"),
		Timestamp: time.Now().Unix(),
	}

	addResp, err := clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: user1Login.Token,
		Data:  data,
	})
	require.NoError(t, err)
	dataID := addResp.Id

	// User2 tries to access user1's data (should fail).
	_, err = clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user2Login.Token,
		Id:    dataID,
	})
	require.Error(t, err) // Should fail - access denied.

	// User2 tries to edit user1's data (should fail).
	_, err = clients.DataClient.EditData(context.Background(), &clientapi.EditDataRequest{
		Token: user2Login.Token,
		Data: &clientapi.Data{
			Id:        dataID,
			Type:      clientapi.DataType_DATA_TYPE_TEXT,
			Payload:   []byte("hacked"),
			Timestamp: time.Now().Unix(),
		},
	})
	require.Error(t, err) // Should fail - access denied.

	// User2 tries to delete user1's data (should fail).
	_, err = clients.DataClient.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: user2Login.Token,
		Id:    dataID,
	})
	require.Error(t, err) // Should fail - access denied.

	// User1 can still access their own data.
	_, err = clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user1Login.Token,
		Id:    dataID,
	})
	require.NoError(t, err) // Should succeed.
}

func TestDataTypes(t *testing.T) {
	clients := testutils.SetupTestClients(t)

	// Setup: register and login user.
	email := testutils.UniqueEmail("testtypes")
	password := "testpassword123"

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

	testCases := []struct {
		name     string
		dataType clientapi.DataType
		metadata interface{}
	}{
		{
			name:     "Login Data",
			dataType: clientapi.DataType_DATA_TYPE_LOGIN,
			metadata: &clientapi.Data_LoginData{
				LoginData: &clientapi.LoginData{
					Username: "testuser",
					Url:      "https://example.com",
					Notes:    "Test login",
				},
			},
		},
		{
			name:     "Card Data",
			dataType: clientapi.DataType_DATA_TYPE_CARD,
			metadata: &clientapi.Data_CardData{
				CardData: &clientapi.CardData{
					CardNumber:     "1234567890123456",
					CardholderName: "Test User",
					ExpiryMonth:    "12",
					ExpiryYear:     "2025",
					Cvv:            "123",
					Notes:          "Test card",
				},
			},
		},
		{
			name:     "Text Data",
			dataType: clientapi.DataType_DATA_TYPE_TEXT,
			metadata: &clientapi.Data_TextData{
				TextData: &clientapi.TextData{
					Title: "Test Note",
					Notes: "This is a test note",
				},
			},
		},
		{
			name:     "Binary Data",
			dataType: clientapi.DataType_DATA_TYPE_BINARY,
			metadata: &clientapi.Data_BinaryData{
				BinaryData: &clientapi.BinaryData{
					Filename:    "test.txt",
					ContentType: "text/plain",
					Size:        1024,
					Notes:       "Test binary file",
				},
			},
		},
		{
			name:     "OTP Data",
			dataType: clientapi.DataType_DATA_TYPE_OTP,
			metadata: &clientapi.Data_OtpData{
				OtpData: &clientapi.OTPData{
					Issuer:    "Test Issuer",
					Account:   "test@example.com",
					Algorithm: "SHA1",
					Digits:    6,
					Period:    30,
					Notes:     "Test OTP",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := &clientapi.Data{
				Type:      tc.dataType,
				Payload:   []byte("test-payload"),
				Timestamp: time.Now().Unix(),
			}

			// Set metadata based on type
			switch md := tc.metadata.(type) {
			case *clientapi.Data_LoginData:
				data.Metadata = md
			case *clientapi.Data_CardData:
				data.Metadata = md
			case *clientapi.Data_TextData:
				data.Metadata = md
			case *clientapi.Data_BinaryData:
				data.Metadata = md
			case *clientapi.Data_OtpData:
				data.Metadata = md
			}

			// Add data
			addResp, err := clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
				Token: token,
				Data:  data,
			})
			require.NoError(t, err)
			require.NotEmpty(t, addResp.Id)

			// Get data and verify
			getResp, err := clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
				Token: token,
				Id:    addResp.Id,
			})
			require.NoError(t, err)
			require.Equal(t, tc.dataType, getResp.Data.Type)
			require.Equal(t, []byte("test-payload"), getResp.Data.Payload)

			// Clean up
			_, err = clients.DataClient.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
				Token: token,
				Id:    addResp.Id,
			})
			require.NoError(t, err)
		})
	}
}
