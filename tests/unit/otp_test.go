package unit

import (
	"testing"
	"time"

	"github.com/AlenaMolokova/gophkeeper/internal/shared/otp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultTOTPConfig(t *testing.T) {
	config := otp.DefaultTOTPConfig()

	assert.Equal(t, "GophKeeper", config.Issuer)
	assert.Equal(t, 20, config.SecretSize)
	assert.Equal(t, 6, config.Digits)
	assert.Equal(t, 30, config.Period)
	assert.Equal(t, "SHA1", config.Algorithm)
}

func TestGenerateSecret(t *testing.T) {
	// Test with default config.
	secret, err := otp.GenerateSecret(nil)
	require.NoError(t, err)
	assert.NotEmpty(t, secret)
	assert.GreaterOrEqual(t, len(secret), 20) // At least 20 characters

	// Test with custom config.
	config := &otp.TOTPConfig{SecretSize: 16}
	secret2, err := otp.GenerateSecret(config)
	require.NoError(t, err)
	assert.NotEmpty(t, secret2)
	assert.GreaterOrEqual(t, len(secret2), 20) // At least 20 characters
}

func TestGenerateAndValidateTOTP(t *testing.T) {
	// Generate a secret.
	secret, err := otp.GenerateSecret(nil)
	require.NoError(t, err)

	// Generate TOTP code.
	code, err := otp.GenerateTOTP(secret, nil)
	require.NoError(t, err)
	assert.Len(t, code, 6) // 6 digits.

	// Validate the code.
	valid := otp.ValidateTOTP(secret, code, nil)
	assert.True(t, valid)

	// Test with invalid code.
	invalid := otp.ValidateTOTP(secret, "123456", nil)
	assert.False(t, invalid)
}

func TestGenerateQRCodeURL(t *testing.T) {
	secret, err := otp.GenerateSecret(nil)
	require.NoError(t, err)

	url, err := otp.GenerateQRCodeURL(secret, "test@example.com", nil)
	require.NoError(t, err)
	assert.Contains(t, url, "otpauth://totp/")
	assert.Contains(t, url, "GophKeeper")
	assert.Contains(t, url, "test@example.com")
}

func TestGenerateBackupCodes(t *testing.T) {
	// Test custom count.
	codes2, err := otp.GenerateBackupCodes(5)
	require.NoError(t, err)
	assert.Len(t, codes2, 5)

	// Test default count (10).
	codes, err := otp.GenerateBackupCodes(10)
	require.NoError(t, err)
	assert.Len(t, codes, 10)

	// Verify code format (8 alphanumeric characters).
	for _, code := range codes2 {
		assert.Len(t, code, 8)
		// Check if it's a valid alphanumeric string.
		for _, char := range code {
			assert.Contains(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", string(char))
		}
	}
}

func TestValidateBackupCode(t *testing.T) {
	codes := []string{"12345678", "87654321", "abcdef12"}

	// Test valid codes.
	assert.True(t, otp.ValidateBackupCode("12345678", codes))
	assert.True(t, otp.ValidateBackupCode("87654321", codes))
	assert.True(t, otp.ValidateBackupCode("abcdef12", codes))

	// Test invalid codes.
	assert.False(t, otp.ValidateBackupCode("11111111", codes))
	assert.False(t, otp.ValidateBackupCode("", codes))
}

func TestTOTPTimeWindow(t *testing.T) {
	secret, err := otp.GenerateSecret(nil)
	require.NoError(t, err)

	// Generate code at current time.
	code, err := otp.GenerateTOTP(secret, nil)
	require.NoError(t, err)

	// Code should be valid immediately.
	assert.True(t, otp.ValidateTOTP(secret, code, nil))

	// Wait for a moment and generate again (should be same code within 30s window).
	time.Sleep(1 * time.Second)
	code2, err := otp.GenerateTOTP(secret, nil)
	require.NoError(t, err)
	assert.Equal(t, code, code2)
}

func TestCustomTOTPConfig(t *testing.T) {
	config := &otp.TOTPConfig{
		SecretSize: 16,
		Digits:     8,
		Period:     60,
	}

	secret, err := otp.GenerateSecret(config)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(secret), 20) // At least 20 characters

	code, err := otp.GenerateTOTP(secret, config)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(code), 6) // At least 6 digits

	// Validate with same config.
	assert.True(t, otp.ValidateTOTP(secret, code, config))
}
