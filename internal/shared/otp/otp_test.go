package otp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultTOTPConfig(t *testing.T) {
	config := DefaultTOTPConfig()
	assert.Equal(t, 30, config.Period)
	assert.Equal(t, 6, config.Digits)
	assert.Equal(t, "SHA1", config.Algorithm)
}

func TestGenerateSecret(t *testing.T) {
	config := DefaultTOTPConfig()
	secret, err := GenerateSecret(config)
	require.NoError(t, err)
	assert.NotEmpty(t, secret)
	assert.GreaterOrEqual(t, len(secret), 20)
}

func TestGenerateAndValidateTOTP(t *testing.T) {
	config := DefaultTOTPConfig()
	secret, err := GenerateSecret(config)
	require.NoError(t, err)

	code, err := GenerateTOTP(secret, config)
	require.NoError(t, err)
	assert.NotEmpty(t, code)
	assert.Len(t, code, config.Digits)

	valid := ValidateTOTP(secret, code, config)
	assert.True(t, valid)
}

func TestGenerateQRCodeURL(t *testing.T) {
	config := DefaultTOTPConfig()
	secret, err := GenerateSecret(config)
	require.NoError(t, err)

	config.Issuer = "TestIssuer"
	config.AccountName = "test@example.com"

	url, err := GenerateQRCodeURL(secret, "test@example.com", config)
	require.NoError(t, err)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "otpauth://totp/")
	assert.Contains(t, url, "TestIssuer")
	assert.Contains(t, url, "test@example.com")
}

func TestGenerateBackupCodes(t *testing.T) {
	codes, err := GenerateBackupCodes(5)
	require.NoError(t, err)
	assert.Len(t, codes, 5)

	for _, code := range codes {
		assert.Len(t, code, 8)
		assert.NotEmpty(t, code)
	}
}

func TestValidateBackupCode(t *testing.T) {
	codes, err := GenerateBackupCodes(3)
	require.NoError(t, err)

	// Test valid code
	valid := ValidateBackupCode(codes[0], codes)
	assert.True(t, valid)

	// Test invalid code
	valid = ValidateBackupCode("INVALID", codes)
	assert.False(t, valid)
}

func TestTOTPTimeWindow(t *testing.T) {
	config := DefaultTOTPConfig()
	secret, err := GenerateSecret(config)
	require.NoError(t, err)

	// Generate code
	code, err := GenerateTOTP(secret, config)
	require.NoError(t, err)
	assert.NotEmpty(t, code)

	// Should be valid immediately
	valid := ValidateTOTP(secret, code, config)
	assert.True(t, valid)

	// Wait for next time window
	time.Sleep(time.Duration(config.Period) * time.Second)

	// Should still be valid in next window
	valid = ValidateTOTP(secret, code, config)
	assert.True(t, valid)
}

func TestCustomTOTPConfig(t *testing.T) {
	config := &TOTPConfig{
		Period:      60,
		Digits:      8,
		Algorithm:   "SHA256",
		Issuer:      "CustomIssuer",
		AccountName: "custom@example.com",
	}

	secret, err := GenerateSecret(config)
	require.NoError(t, err)

	code, err := GenerateTOTP(secret, config)
	require.NoError(t, err)
	assert.NotEmpty(t, code)

	valid := ValidateTOTP(secret, code, config)
	assert.True(t, valid)
}
