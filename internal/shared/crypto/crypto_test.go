package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptData(t *testing.T) {
	key := make([]byte, 32) // 32 bytes for AES-256
	for i := range key {
		key[i] = byte(i)
	}

	testData := []byte("test data")
	encrypted, err := EncryptData(key, testData)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, testData, encrypted)
}

func TestDecryptData(t *testing.T) {
	key := make([]byte, 32) // 32 bytes for AES-256
	for i := range key {
		key[i] = byte(i)
	}

	testData := []byte("test data")
	encrypted, err := EncryptData(key, testData)
	require.NoError(t, err)

	decrypted, err := DecryptData(key, encrypted)
	require.NoError(t, err)
	assert.Equal(t, testData, decrypted)
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := make([]byte, 32) // 32 bytes for AES-256
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
			testData := []byte(tc)
			encrypted, err := EncryptData(key, testData)
			require.NoError(t, err)

			decrypted, err := DecryptData(key, encrypted)
			require.NoError(t, err)
			assert.Equal(t, testData, decrypted)
		})
	}
}
