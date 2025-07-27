package integration

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/AlenaMolokova/gophkeeper/tests/integration/testutils"
	"github.com/stretchr/testify/require"
)

func uniqueEmail(base string) string {
	return base + strconv.Itoa(rand.Intn(1_000_000)) + "@example.com"
}

func TestSyncBetweenClients(t *testing.T) {
	clients := testutils.SetupTestClients(t)

	// Setup: register and login user.
	email := uniqueEmail("testsync")
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

	// Add data from "client 1".
	data1 := &clientapi.Data{
		Type:      clientapi.DataType_DATA_TYPE_LOGIN,
		Payload:   []byte("client1-data"),
		Metadata:  &clientapi.Data_LoginData{LoginData: &clientapi.LoginData{Notes: "client1"}},
		Timestamp: time.Now().Unix(),
	}

	addResp1, err := clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: token,
		Data:  data1,
	})
	require.NoError(t, err)

	// Simulate "client 2" adding different data.
	data2 := &clientapi.Data{
		Type:      clientapi.DataType_DATA_TYPE_TEXT,
		Payload:   []byte("client2-data"),
		Metadata:  &clientapi.Data_TextData{TextData: &clientapi.TextData{Notes: "client2"}},
		Timestamp: time.Now().Unix(),
	}

	addResp2, err := clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: token,
		Data:  data2,
	})
	require.NoError(t, err)

	// Verify both data items exist.
	getResp1, err := clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: token,
		Id:    addResp1.Id,
	})
	require.NoError(t, err)
	require.Equal(t, "client1-data", string(getResp1.Data.Payload))

	getResp2, err := clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: token,
		Id:    addResp2.Id,
	})
	require.NoError(t, err)
	require.Equal(t, "client2-data", string(getResp2.Data.Payload))

	// Clean up.
	_, err = clients.DataClient.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: token,
		Id:    addResp1.Id,
	})
	require.NoError(t, err)

	_, err = clients.DataClient.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: token,
		Id:    addResp2.Id,
	})
	require.NoError(t, err)
}

func TestMultiUserSync(t *testing.T) {
	clients := testutils.SetupTestClients(t)

	// Setup: create two users.
	user1Email := uniqueEmail("syncuser1")
	user2Email := uniqueEmail("syncuser2")
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

	// User1 adds data.
	user1Data := &clientapi.Data{
		Type:      clientapi.DataType_DATA_TYPE_TEXT,
		Payload:   []byte("user1-data"),
		Metadata:  &clientapi.Data_TextData{TextData: &clientapi.TextData{Notes: "user1"}},
		Timestamp: time.Now().Unix(),
	}

	addResp1, err := clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: user1Login.Token,
		Data:  user1Data,
	})
	require.NoError(t, err)

	// User2 adds different data.
	user2Data := &clientapi.Data{
		Type:      clientapi.DataType_DATA_TYPE_TEXT,
		Payload:   []byte("user2-data"),
		Metadata:  &clientapi.Data_TextData{TextData: &clientapi.TextData{Notes: "user2"}},
		Timestamp: time.Now().Unix(),
	}

	addResp2, err := clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: user2Login.Token,
		Data:  user2Data,
	})
	require.NoError(t, err)

	// Verify each user can only access their own data.
	getResp1, err := clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user1Login.Token,
		Id:    addResp1.Id,
	})
	require.NoError(t, err)
	require.Equal(t, "user1-data", string(getResp1.Data.Payload))

	getResp2, err := clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user2Login.Token,
		Id:    addResp2.Id,
	})
	require.NoError(t, err)
	require.Equal(t, "user2-data", string(getResp2.Data.Payload))

	// Verify users cannot access each other's data.
	_, err = clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user1Login.Token,
		Id:    addResp2.Id,
	})
	require.Error(t, err) // Should fail - access denied.

	_, err = clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user2Login.Token,
		Id:    addResp1.Id,
	})
	require.Error(t, err) // Should fail - access denied.

	// Clean up.
	_, err = clients.DataClient.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: user1Login.Token,
		Id:    addResp1.Id,
	})
	require.NoError(t, err)

	_, err = clients.DataClient.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: user2Login.Token,
		Id:    addResp2.Id,
	})
	require.NoError(t, err)
}

func TestDataConsistency(t *testing.T) {
	clients := testutils.SetupTestClients(t)

	// Setup: register and login user.
	email := uniqueEmail("testconsistency")
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

	// Add data and verify consistency.
	data := &clientapi.Data{
		Type:      clientapi.DataType_DATA_TYPE_TEXT,
		Payload:   []byte("consistent-data"),
		Metadata:  &clientapi.Data_TextData{TextData: &clientapi.TextData{Title: "Test"}},
		Timestamp: time.Now().Unix(),
	}

	addResp, err := clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: token,
		Data:  data,
	})
	require.NoError(t, err)

	// Read the data multiple times to ensure consistency.
	for i := 0; i < 5; i++ {
		getResp, err := clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
			Token: token,
			Id:    addResp.Id,
		})
		require.NoError(t, err)
		require.Equal(t, "consistent-data", string(getResp.Data.Payload))
		require.Equal(t, clientapi.DataType_DATA_TYPE_TEXT, getResp.Data.Type)
	}

	// Clean up.
	_, err = clients.DataClient.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: token,
		Id:    addResp.Id,
	})
	require.NoError(t, err)
}
