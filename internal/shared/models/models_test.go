package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user := User{
		ID:    "user-123",
		Email: "test@example.com",
		Hash:  "hashed-password",
	}

	assert.Equal(t, "user-123", user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashed-password", user.Hash)
}

func TestData(t *testing.T) {
	now := time.Now().Unix()
	data := Data{
		ID:        "data-123",
		Type:      "login",
		Payload:   []byte("encrypted-payload"),
		Metadata:  map[string]string{"key": "value"},
		Timestamp: now,
	}

	assert.Equal(t, "data-123", data.ID)
	assert.Equal(t, "login", data.Type)
	assert.Equal(t, []byte("encrypted-payload"), data.Payload)
	assert.Equal(t, map[string]string{"key": "value"}, data.Metadata)
	assert.Equal(t, now, data.Timestamp)
}

func TestDataWithEmptyMetadata(t *testing.T) {
	data := Data{
		ID:        "data-456",
		Type:      "text",
		Payload:   []byte("test-data"),
		Metadata:  make(map[string]string),
		Timestamp: time.Now().Unix(),
	}

	assert.Equal(t, "data-456", data.ID)
	assert.Equal(t, "text", data.Type)
	assert.Equal(t, []byte("test-data"), data.Payload)
	assert.Empty(t, data.Metadata)
	assert.Greater(t, data.Timestamp, int64(0))
}

func TestUserWithEmptyFields(t *testing.T) {
	user := User{}

	assert.Empty(t, user.ID)
	assert.Empty(t, user.Email)
	assert.Empty(t, user.Hash)
}

func TestDataWithEmptyFields(t *testing.T) {
	data := Data{}

	assert.Empty(t, data.ID)
	assert.Empty(t, data.Type)
	assert.Nil(t, data.Payload)
	assert.Nil(t, data.Metadata)
	assert.Equal(t, int64(0), data.Timestamp)
}
