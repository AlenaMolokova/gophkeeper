// Package crypto provides client-side cryptographic operations for the GophKeeper application.
// It handles key generation, storage, and provides high-level encryption/decryption functions.
package crypto

import (
	"crypto/rand"
	"errors"
	"os"
	"path/filepath"

	sharedcrypto "github.com/AlenaMolokova/gophkeeper/internal/shared/crypto"
)

const keyFileName = "gophkeeper.key"

// GenerateKey creates a new 32-byte AES-256 encryption key.
// Returns the generated key as a byte slice, or an error if key generation fails.
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// SaveKey saves the encryption key to a file in the user's home directory.
// The key file is created with restrictive permissions (0600) for security.
// Returns an error if the key cannot be saved.
func SaveKey(key []byte) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, keyFileName)
	return os.WriteFile(path, key, 0o600)
}

// LoadKey loads the encryption key from the file in the user's home directory.
// Returns the loaded key as a byte slice, or an error if the key cannot be loaded.
// Validates that the loaded key is exactly 32 bytes long.
func LoadKey() ([]byte, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, keyFileName)
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, errors.New("ключ поврежден или неверного размера")
	}
	return key, nil
}

// EncryptData encrypts data using the user's stored encryption key.
// Automatically loads the key and calls the shared encryption function.
// Returns the encrypted data as a byte slice, or an error if encryption fails.
func EncryptData(plaintext []byte) ([]byte, error) {
	key, err := LoadKey()
	if err != nil {
		return nil, err
	}
	return sharedcrypto.EncryptData(key, plaintext)
}

// DecryptData decrypts data using the user's stored encryption key.
// Automatically loads the key and calls the shared decryption function.
// Returns the decrypted data as a byte slice, or an error if decryption fails.
func DecryptData(ciphertext []byte) ([]byte, error) {
	key, err := LoadKey()
	if err != nil {
		return nil, err
	}
	return sharedcrypto.DecryptData(key, ciphertext)
}
