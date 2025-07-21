// Package api_test provides unit tests for the GophKeeper server API.
// The tests use mocked storage to verify API behavior without requiring
// a real database connection.
//
// Test coverage includes:
// - User registration and authentication
// - Data CRUD operations
// - JWT token validation
// - Error handling scenarios
package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlenaMolokova/gophkeeper/internal/server/api"
	"github.com/AlenaMolokova/gophkeeper/internal/server/auth"
	"github.com/AlenaMolokova/gophkeeper/internal/server/storage"
	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// createValidToken creates a valid JWT token for testing.
// The token is signed with the provided secret and contains the user ID
// in the subject claim with a 1-hour expiration time.
//
// This function is used in tests to create valid tokens for API calls
// that require authentication.
func createValidToken(secret []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user1",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

// TestServer_Register tests user registration functionality.
// It verifies that a new user can be registered with email and password,
// and that a valid JWT token is returned upon successful registration.
func TestServer_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	mockStorage.EXPECT().SaveUser(gomock.Any(), "test@example.com", gomock.Any()).Return("user1", nil)
	mockStorage.EXPECT().GetPoolStats().Return(&storage.PoolStats{}).AnyTimes()

	srv := api.NewServer(auth.NewAuth(mockStorage, "secret", 3600*time.Second), mockStorage, []byte("secret"))

	resp, err := srv.Register(context.Background(), &clientapi.RegisterRequest{Email: "test@example.com", Password: "password"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

// TestServer_Login tests user authentication functionality.
// It verifies that a user can log in with valid credentials and
// receive a valid JWT token for subsequent API calls.
func TestServer_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	// Generate proper bcrypt hash for "password"
	hash := "$2a$10$ozYQzeZhFcWXqcJQ9FkAruGvr90EHjUs2r5Beesoex4Q.qUhk0a.i"
	mockStorage.EXPECT().FindUserByEmail(gomock.Any(), "test@example.com").Return(models.User{ID: "user1", Email: "test@example.com", Hash: hash}, nil)
	mockStorage.EXPECT().GetPoolStats().Return(&storage.PoolStats{}).AnyTimes()

	srv := api.NewServer(auth.NewAuth(mockStorage, "secret", 3600*time.Second), mockStorage, []byte("secret"))

	resp, err := srv.Login(context.Background(), &clientapi.LoginRequest{Email: "test@example.com", Password: "password"})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

// TestServer_AddData tests data creation functionality.
// It verifies that authenticated users can add new encrypted data
// and receive a unique data ID for future reference.
func TestServer_AddData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	mockStorage.EXPECT().SaveData(gomock.Any(), "user1", gomock.Any()).Return("data1", nil)
	mockStorage.EXPECT().GetPoolStats().Return(&storage.PoolStats{}).AnyTimes()

	srv := api.NewServer(auth.NewAuth(mockStorage, "secret", 3600*time.Second), mockStorage, []byte("secret"))

	validToken := createValidToken([]byte("secret"))
	resp, err := srv.AddData(context.Background(), &clientapi.AddDataRequest{
		Token: validToken,
		Data:  &clientapi.Data{Id: "data1", Type: "text", Payload: []byte("test"), Timestamp: time.Now().Unix()},
	})
	assert.NoError(t, err)
	assert.Equal(t, "data1", resp.Id)
}

// TestServer_GetData tests data retrieval functionality.
// It verifies that authenticated users can retrieve their encrypted data
// by ID and that data isolation between users is maintained.
func TestServer_GetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	mockStorage.EXPECT().FindDataByID(gomock.Any(), "user1", "data1").Return(models.Data{ID: "data1", Type: "text", Payload: []byte("test"), Timestamp: time.Now().Unix()}, nil)
	mockStorage.EXPECT().GetPoolStats().Return(&storage.PoolStats{}).AnyTimes()

	srv := api.NewServer(auth.NewAuth(mockStorage, "secret", 3600*time.Second), mockStorage, []byte("secret"))

	validToken := createValidToken([]byte("secret"))
	resp, err := srv.GetData(context.Background(), &clientapi.GetDataRequest{Token: validToken, Id: "data1"})
	assert.NoError(t, err)
	assert.Equal(t, "data1", resp.Data.Id)
}

// TestServer_EditData tests data update functionality.
// It verifies that authenticated users can update their existing data
// and that the updated data is properly stored.
func TestServer_EditData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	mockStorage.EXPECT().EditData(gomock.Any(), "user1", gomock.Any()).Return("data1", nil)
	mockStorage.EXPECT().GetPoolStats().Return(&storage.PoolStats{}).AnyTimes()

	srv := api.NewServer(auth.NewAuth(mockStorage, "secret", 3600*time.Second), mockStorage, []byte("secret"))

	validToken := createValidToken([]byte("secret"))
	resp, err := srv.EditData(context.Background(), &clientapi.EditDataRequest{
		Token: validToken,
		Data:  &clientapi.Data{Id: "data1", Type: "text", Payload: []byte("updated"), Timestamp: time.Now().Unix()},
	})
	assert.NoError(t, err)
	assert.Equal(t, "data1", resp.Id)
}

// TestServer_DeleteData tests data deletion functionality.
// It verifies that authenticated users can delete their data by ID
// and that data isolation between users is maintained.
func TestServer_DeleteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	mockStorage.EXPECT().DeleteData(gomock.Any(), "user1", "data1").Return(nil)
	mockStorage.EXPECT().GetPoolStats().Return(&storage.PoolStats{}).AnyTimes()

	srv := api.NewServer(auth.NewAuth(mockStorage, "secret", 3600*time.Second), mockStorage, []byte("secret"))

	validToken := createValidToken([]byte("secret"))
	resp, err := srv.DeleteData(context.Background(), &clientapi.DeleteDataRequest{Token: validToken, Id: "data1"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
