// Package crypto provides encryption and decryption functionality for the GophKeeper application.
// It implements AES-256 encryption for securing sensitive data on the client side.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// EncryptData encrypts plaintext data using AES-256 encryption.
// The key must be exactly 32 bytes long for AES-256.
// Returns the encrypted data as a byte slice, or an error if encryption fails.
func EncryptData(key, plaintext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("ключ должен быть 32 байта (AES-256)")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

// DecryptData decrypts ciphertext data using AES-256 decryption.
// The key must be exactly 32 bytes long for AES-256.
// Returns the decrypted data as a byte slice, or an error if decryption fails.
func DecryptData(key, ciphertext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("ключ должен быть 32 байта (AES-256)")
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext слишком короткий")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
