// Package integration provides parallel integration tests using Testcontainers.
// These tests demonstrate how to run multiple test scenarios in parallel
// with complete isolation through separate containers.
package integration

import (
	"context"
	"sync"
	"testing"
	"time"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/AlenaMolokova/gophkeeper/tests/integration/testutils"
	"github.com/stretchr/testify/assert"
)

// TestParallelUserRegistration runs multiple user registration tests in parallel.
// Each test uses its own isolated PostgreSQL container and TLS certificates.
func TestParallelUserRegistration(t *testing.T) {
	// Run multiple registration tests in parallel
	t.Run("Parallel", func(t *testing.T) {
		const numTests = 5
		var wg sync.WaitGroup
		results := make([]bool, numTests)

		for i := 0; i < numTests; i++ {
			wg.Add(1)
			go func(testIndex int) {
				defer wg.Done()
				results[testIndex] = runUserRegistrationTest(t, testIndex)
			}(i)
		}

		wg.Wait()

		// Verify all tests passed
		for i, result := range results {
			assert.True(t, result, "Test %d should have passed", i)
		}
	})
}

// TestParallelDataOperations runs multiple data operations in parallel.
// Demonstrates concurrent data access with proper isolation.
func TestParallelDataOperations(t *testing.T) {
	t.Run("Parallel", func(t *testing.T) {
		const numTests = 3
		var wg sync.WaitGroup
		results := make([]bool, numTests)

		for i := 0; i < numTests; i++ {
			wg.Add(1)
			go func(testIndex int) {
				defer wg.Done()
				results[testIndex] = runDataOperationsTest(t, testIndex)
			}(i)
		}

		wg.Wait()

		// Verify all tests passed
		for i, result := range results {
			assert.True(t, result, "Test %d should have passed", i)
		}
	})
}

// TestConcurrentUserAccess tests multiple users accessing the system simultaneously.
// Each user gets their own isolated environment.
func TestConcurrentUserAccess(t *testing.T) {
	t.Run("Concurrent", func(t *testing.T) {
		const numUsers = 4
		var wg sync.WaitGroup
		results := make([]bool, numUsers)

		for i := 0; i < numUsers; i++ {
			wg.Add(1)
			go func(userIndex int) {
				defer wg.Done()
				results[userIndex] = runConcurrentUserTest(t, userIndex)
			}(i)
		}

		wg.Wait()

		// Verify all user tests passed
		for i, result := range results {
			assert.True(t, result, "User test %d should have passed", i)
		}
	})
}

// runUserRegistrationTest executes a single user registration test in isolation.
func runUserRegistrationTest(t *testing.T, testIndex int) bool {
	// Setup isolated test environment with Testcontainers
	env := testutils.SetupTestEnvironment(t)
	defer env.Cleanup()

	// Setup test server with containerized database
	testutils.SetupTestServer(t, env)

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	// Setup test clients with dynamic certificates
	clients := testutils.SetupTestClients(t, env)

	// Generate unique email for this test
	email := testutils.UniqueEmail("paralleluser")
	password := "testpassword123"

	// Test registration
	resp, err := clients.UserClient.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Logf("Test %d: Registration failed: %v", testIndex, err)
		return false
	}

	if resp.Token == "" {
		t.Logf("Test %d: No token returned", testIndex)
		return false
	}

	// Test login with the same credentials
	loginResp, err := clients.UserClient.Login(context.Background(), &clientapi.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Logf("Test %d: Login failed: %v", testIndex, err)
		return false
	}

	if loginResp.Token == "" {
		t.Logf("Test %d: No login token returned", testIndex)
		return false
	}

	return true
}

// runDataOperationsTest executes data CRUD operations in isolation.
func runDataOperationsTest(t *testing.T, testIndex int) bool {
	// Setup isolated test environment with Testcontainers
	env := testutils.SetupTestEnvironment(t)
	defer env.Cleanup()

	// Setup test server with containerized database
	testutils.SetupTestServer(t, env)

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	// Setup test clients with dynamic certificates
	clients := testutils.SetupTestClients(t, env)

	// Register and login user
	email := testutils.UniqueEmail("paralleluser")
	password := "testpassword123"

	_, err := clients.UserClient.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Logf("Test %d: Registration failed: %v", testIndex, err)
		return false
	}

	loginResp, err := clients.UserClient.Login(context.Background(), &clientapi.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Logf("Test %d: Login failed: %v", testIndex, err)
		return false
	}

	token := loginResp.Token

	// Add data
	addResp, err := clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: token,
		Data: &clientapi.Data{
			Type:      clientapi.DataType_DATA_TYPE_TEXT,
			Payload:   []byte("test data"),
			Timestamp: time.Now().Unix(),
		},
	})
	if err != nil {
		t.Logf("Test %d: Add data failed: %v", testIndex, err)
		return false
	}

	// Get data
	_, err = clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
		Token: token,
		Id:    addResp.Id,
	})
	if err != nil {
		t.Logf("Test %d: Get data failed: %v", testIndex, err)
		return false
	}

	// Delete data
	_, err = clients.DataClient.DeleteData(context.Background(), &clientapi.DeleteDataRequest{
		Token: token,
		Id:    addResp.Id,
	})
	if err != nil {
		t.Logf("Test %d: Delete data failed: %v", testIndex, err)
		return false
	}

	return true
}

// runConcurrentUserTest tests a single user's operations in isolation.
func runConcurrentUserTest(t *testing.T, userIndex int) bool {
	// Setup isolated test environment with Testcontainers
	env := testutils.SetupTestEnvironment(t)
	defer env.Cleanup()

	// Setup test server with containerized database
	testutils.SetupTestServer(t, env)

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	// Setup test clients with dynamic certificates
	clients := testutils.SetupTestClients(t, env)

	// Register and login user
	email := testutils.UniqueEmail("concurrentuser")
	password := "testpassword123"

	_, err := clients.UserClient.Register(context.Background(), &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Logf("User %d: Registration failed: %v", userIndex, err)
		return false
	}

	loginResp, err := clients.UserClient.Login(context.Background(), &clientapi.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Logf("User %d: Login failed: %v", userIndex, err)
		return false
	}

	token := loginResp.Token

	// Perform multiple operations
	for i := 0; i < 3; i++ {
		// Add data
		addResp, err := clients.DataClient.AddData(context.Background(), &clientapi.AddDataRequest{
			Token: token,
			Data: &clientapi.Data{
				Type:      clientapi.DataType_DATA_TYPE_TEXT,
				Payload:   []byte("user data"),
				Timestamp: time.Now().Unix(),
			},
		})
		if err != nil {
			t.Logf("User %d: Add data %d failed: %v", userIndex, i, err)
			return false
		}

		// Verify data
		_, err = clients.DataClient.GetData(context.Background(), &clientapi.GetDataRequest{
			Token: token,
			Id:    addResp.Id,
		})
		if err != nil {
			t.Logf("User %d: Get data %d failed: %v", userIndex, i, err)
			return false
		}
	}

	return true
}
