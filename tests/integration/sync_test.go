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
	client := testutils.SetupTestClient(t)

	// Setup: register and login user.
	email := uniqueEmail("testsync")
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

	// Add data from "client 1".
	data1 := &clientapi.Data{
		Type:      "login",
		Payload:   []byte("client1-data"),
		Metadata:  map[string]string{"source": "client1"},
		Timestamp: time.Now().Unix(),
	}

	addResp1, err := client.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: token,
		Data:  data1,
	})
	require.NoError(t, err)

	// Simulate "client 2" adding different data.
	data2 := &clientapi.Data{
		Type:      "text",
		Payload:   []byte("client2-data"),
		Metadata:  map[string]string{"source": "client2"},
		Timestamp: time.Now().Unix(),
	}

	addResp2, err := client.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: token,
		Data:  data2,
	})
	require.NoError(t, err)

	// Verify both data items exist.
	getResp1, err := client.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: token,
		Id:    addResp1.Id,
	})
	require.NoError(t, err)
	require.Equal(t, "client1-data", string(getResp1.Data.Payload))

	getResp2, err := client.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: token,
		Id:    addResp2.Id,
	})
	require.NoError(t, err)
	require.Equal(t, "client2-data", string(getResp2.Data.Payload))

	// Clean up.
	_, err = client.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: token,
		Id:    addResp1.Id,
	})
	require.NoError(t, err)

	_, err = client.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: token,
		Id:    addResp2.Id,
	})
	require.NoError(t, err)
}

func TestMultiUserSync(t *testing.T) {
	client := testutils.SetupTestClient(t)

	// Setup: create two users.
	user1Email := uniqueEmail("syncuser1")
	user2Email := uniqueEmail("syncuser2")
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

	// User1 adds data.
	data1 := &clientapi.Data{
		Type:      "user1-data",
		Payload:   []byte("user1-secret"),
		Metadata:  map[string]string{"owner": "user1"},
		Timestamp: time.Now().Unix(),
	}

	addResp1, err := client.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: user1Login.Token,
		Data:  data1,
	})
	require.NoError(t, err)

	// User2 adds different data.
	data2 := &clientapi.Data{
		Type:      "user2-data",
		Payload:   []byte("user2-secret"),
		Metadata:  map[string]string{"owner": "user2"},
		Timestamp: time.Now().Unix(),
	}

	addResp2, err := client.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: user2Login.Token,
		Data:  data2,
	})
	require.NoError(t, err)

	// Verify users can only access their own data.
	getResp1, err := client.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user1Login.Token,
		Id:    addResp1.Id,
	})
	require.NoError(t, err)
	require.Equal(t, "user1-secret", string(getResp1.Data.Payload))

	getResp2, err := client.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user2Login.Token,
		Id:    addResp2.Id,
	})
	require.NoError(t, err)
	require.Equal(t, "user2-secret", string(getResp2.Data.Payload))

	// Verify users cannot access each other's data.
	_, err = client.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user1Login.Token,
		Id:    addResp2.Id,
	})
	require.Error(t, err) // User1 cannot access User2's data.

	_, err = client.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: user2Login.Token,
		Id:    addResp1.Id,
	})
	require.Error(t, err) // User2 cannot access User1's data.

	// Clean up.
	_, err = client.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: user1Login.Token,
		Id:    addResp1.Id,
	})
	require.NoError(t, err)

	_, err = client.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: user2Login.Token,
		Id:    addResp2.Id,
	})
	require.NoError(t, err)
}

func TestDataConsistency(t *testing.T) {
	client := testutils.SetupTestClient(t)

	// Setup: register and login user.
	email := uniqueEmail("testconsistency")
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

	// Test data consistency across multiple operations.
	testData := []byte("consistent-data")
	data := &clientapi.Data{
		Type:      "consistency-test",
		Payload:   testData,
		Metadata:  map[string]string{"test": "consistency"},
		Timestamp: time.Now().Unix(),
	}

	// Add data.
	addResp, err := client.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: token,
		Data:  data,
	})
	require.NoError(t, err)
	dataID := addResp.Id

	// Read data multiple times to ensure consistency.
	for i := 0; i < 5; i++ {
		getResp, err := client.GetData(context.Background(), &clientapi.GetDataRequest{
			Token: token,
			Id:    dataID,
		})
		require.NoError(t, err)
		require.Equal(t, testData, getResp.Data.Payload)
		require.Equal(t, "consistency-test", getResp.Data.Type)
		require.Equal(t, "consistency", getResp.Data.Metadata["test"])
	}

	// Update data.
	updatedData := []byte("updated-consistent-data")
	updateData := &clientapi.Data{
		Id:        dataID,
		Type:      "consistency-test",
		Payload:   updatedData,
		Metadata:  map[string]string{"test": "consistency", "updated": "true"},
		Timestamp: time.Now().Unix(),
	}

	_, err = client.EditData(context.Background(), &clientapi.EditDataRequest{
		Token: token,
		Data:  updateData,
	})
	require.NoError(t, err)

	// Read updated data multiple times to ensure consistency.
	for i := 0; i < 5; i++ {
		getResp, err := client.GetData(context.Background(), &clientapi.GetDataRequest{
			Token: token,
			Id:    dataID,
		})
		require.NoError(t, err)
		require.Equal(t, updatedData, getResp.Data.Payload)
		require.Equal(t, "consistency-test", getResp.Data.Type)
		require.Equal(t, "consistency", getResp.Data.Metadata["test"])
		require.Equal(t, "true", getResp.Data.Metadata["updated"])
	}

	// Clean up.
	_, err = client.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: token,
		Id:    dataID,
	})
	require.NoError(t, err)
}
