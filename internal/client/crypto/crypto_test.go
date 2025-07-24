package crypto

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateKey(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)
	assert.Len(t, key, 32)

	// Generate another key to ensure they're different
	key2, err := GenerateKey()
	require.NoError(t, err)
	assert.Len(t, key2, 32)
	assert.NotEqual(t, key, key2)
}

func TestSaveKeyAndLoadKey(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)

	// Also set USERPROFILE for Windows
	originalUserProfile := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", originalUserProfile)
	os.Setenv("USERPROFILE", tempDir)

	// Generate and save key
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i)
	}

	err := SaveKey(testKey)
	require.NoError(t, err)

	// Check that file was created
	keyPath := filepath.Join(tempDir, keyFileName)
	_, err = os.Stat(keyPath)
	require.NoError(t, err)

	// Load the key
	loadedKey, err := LoadKey()
	require.NoError(t, err)
	assert.Equal(t, testKey, loadedKey)
}

func TestLoadKeyWithInvalidSize(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)

	// Also set USERPROFILE for Windows
	originalUserProfile := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", originalUserProfile)
	os.Setenv("USERPROFILE", tempDir)

	// Create a key file with wrong size
	keyPath := filepath.Join(tempDir, keyFileName)
	err := os.WriteFile(keyPath, []byte("short"), 0o600)
	require.NoError(t, err)

	// Try to load the key
	_, err = LoadKey()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ключ поврежден или неверного размера")
}

func TestLoadKeyFileNotFound(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)

	// Also set USERPROFILE for Windows
	originalUserProfile := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", originalUserProfile)
	os.Setenv("USERPROFILE", tempDir)

	// Try to load non-existent key
	_, err := LoadKey()
	assert.Error(t, err)
}

func TestEncryptDataAndDecryptData(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)

	// Also set USERPROFILE for Windows
	originalUserProfile := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", originalUserProfile)
	os.Setenv("USERPROFILE", tempDir)

	// Generate and save a key
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i)
	}
	err := SaveKey(testKey)
	require.NoError(t, err)

	// Test data
	testData := []byte("test data for encryption")

	// Encrypt data
	encrypted, err := EncryptData(testData)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, testData, encrypted)

	// Decrypt data
	decrypted, err := DecryptData(encrypted)
	require.NoError(t, err)
	assert.Equal(t, testData, decrypted)
}

func TestEncryptDataWithoutKey(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)

	// Also set USERPROFILE for Windows
	originalUserProfile := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", originalUserProfile)
	os.Setenv("USERPROFILE", tempDir)

	// Try to encrypt without saving a key first
	testData := []byte("test data")
	_, err := EncryptData(testData)
	assert.Error(t, err)
}

func TestDecryptDataWithoutKey(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)

	// Also set USERPROFILE for Windows
	originalUserProfile := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", originalUserProfile)
	os.Setenv("USERPROFILE", tempDir)

	// Try to decrypt without saving a key first
	testData := []byte("test data")
	_, err := DecryptData(testData)
	assert.Error(t, err)
}
