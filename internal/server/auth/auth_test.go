package auth

import (
	"context"
	"testing"
	"time"

	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockStorage is a mock implementation of auth.Storage for testing.
type MockStorage struct {
	users map[string]models.User
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		users: make(map[string]models.User),
	}
}

func (m *MockStorage) SaveUser(ctx context.Context, email, hash string) (string, error) {
	userID := "user-" + email
	m.users[email] = models.User{
		ID:    userID,
		Email: email,
		Hash:  hash,
	}
	return userID, nil
}

func (m *MockStorage) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, exists := m.users[email]
	if !exists {
		return models.User{}, assert.AnError
	}
	return user, nil
}

func TestNewAuth(t *testing.T) {
	mockStorage := NewMockStorage()
	jwtSecret := "test-secret"
	tokenTTL := time.Hour

	authService := NewAuth(mockStorage, jwtSecret, tokenTTL)
	require.NotNil(t, authService)
}

func TestRegister(t *testing.T) {
	mockStorage := NewMockStorage()
	authService := NewAuth(mockStorage, "test-secret", time.Hour)

	ctx := context.Background()
	token, err := authService.Register(ctx, "test@example.com", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Check that user was saved
	user, err := mockStorage.FindUserByEmail(ctx, "test@example.com")
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEmpty(t, user.Hash)
}

func TestLogin(t *testing.T) {
	mockStorage := NewMockStorage()
	authService := NewAuth(mockStorage, "test-secret", time.Hour)

	ctx := context.Background()

	// First register a user
	_, err := authService.Register(ctx, "test@example.com", "password123")
	require.NoError(t, err)

	// Then login
	token, err := authService.Login(ctx, "test@example.com", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLoginWithWrongPassword(t *testing.T) {
	mockStorage := NewMockStorage()
	authService := NewAuth(mockStorage, "test-secret", time.Hour)

	ctx := context.Background()

	// First register a user
	_, err := authService.Register(ctx, "test@example.com", "password123")
	require.NoError(t, err)

	// Then try to login with wrong password
	_, err = authService.Login(ctx, "test@example.com", "wrongpassword")
	assert.Error(t, err)
}

func TestLoginWithNonExistentUser(t *testing.T) {
	mockStorage := NewMockStorage()
	authService := NewAuth(mockStorage, "test-secret", time.Hour)

	ctx := context.Background()
	_, err := authService.Login(ctx, "nonexistent@example.com", "password123")
	assert.Error(t, err)
}
