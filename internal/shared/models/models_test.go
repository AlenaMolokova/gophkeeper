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
	loginData := &LoginData{
		Username: "testuser",
		Url:      "https://example.com",
		Notes:    "Test login",
	}

	data := Data{
		ID:        "data-123",
		Type:      DataTypeLogin,
		Payload:   []byte("encrypted-payload"),
		Metadata:  loginData,
		Timestamp: now,
	}

	assert.Equal(t, "data-123", data.ID)
	assert.Equal(t, DataTypeLogin, data.Type)
	assert.Equal(t, []byte("encrypted-payload"), data.Payload)
	assert.Equal(t, loginData, data.Metadata)
	assert.Equal(t, now, data.Timestamp)
}

func TestDataWithEmptyMetadata(t *testing.T) {
	textData := &TextData{
		Title: "Test Title",
		Notes: "Test notes",
	}

	data := Data{
		ID:        "data-456",
		Type:      DataTypeText,
		Payload:   []byte("test-data"),
		Metadata:  textData,
		Timestamp: time.Now().Unix(),
	}

	assert.Equal(t, "data-456", data.ID)
	assert.Equal(t, DataTypeText, data.Type)
	assert.Equal(t, []byte("test-data"), data.Payload)
	assert.Equal(t, textData, data.Metadata)
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
	assert.Equal(t, DataTypeUnspecified, data.Type)
	assert.Nil(t, data.Payload)
	assert.Nil(t, data.Metadata)
	assert.Equal(t, int64(0), data.Timestamp)
}

func TestDataTypes(t *testing.T) {
	// Test all data types
	types := []DataType{
		DataTypeLogin,
		DataTypeText,
		DataTypeBinary,
		DataTypeCard,
		DataTypeOTP,
	}

	for _, dataType := range types {
		data := Data{
			ID:   "test-id",
			Type: dataType,
		}
		assert.Equal(t, dataType, data.Type)
	}
}
