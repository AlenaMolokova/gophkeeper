package unit

import (
	"testing"

	"github.com/AlenaMolokova/gophkeeper/internal/shared/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptData(t *testing.T) {
	// Test with valid key and data.
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	plaintext := []byte("test data")

	encrypted, err := crypto.EncryptData(key, plaintext)
	require.NoError(t, err)
	assert.NotNil(t, encrypted)
	assert.NotEqual(t, plaintext, encrypted)
	assert.Greater(t, len(encrypted), len(plaintext))

	// Test with invalid key length.
	invalidKey := []byte("short")
	_, err = crypto.EncryptData(invalidKey, plaintext)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 байта")
}

func TestDecryptData(t *testing.T) {
	// Test with valid key and data.
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	plaintext := []byte("test data")

	encrypted, err := crypto.EncryptData(key, plaintext)
	require.NoError(t, err)

	decrypted, err := crypto.DecryptData(key, encrypted)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)

	// Test with invalid key length.
	invalidKey := []byte("short")
	_, err = crypto.DecryptData(invalidKey, encrypted)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 байта")

	// Test with too short ciphertext.
	_, err = crypto.DecryptData(key, []byte("short"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "слишком короткий")
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	testCases := []string{
		"",
		"hello world",
		"привет мир",
		"special chars: !@#$%^&*()",
		"very long string with many characters to test encryption and decryption functionality",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			plaintext := []byte(tc)
			encrypted, err := crypto.EncryptData(key, plaintext)
			require.NoError(t, err)

			decrypted, err := crypto.DecryptData(key, encrypted)
			require.NoError(t, err)
			assert.Equal(t, plaintext, decrypted)
		})
	}
}
