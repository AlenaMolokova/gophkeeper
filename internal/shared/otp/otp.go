package otp

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// TOTPConfig holds configuration for TOTP generation and validation.
type TOTPConfig struct {
	// Issuer is the name of the service (e.g., "GophKeeper").
	Issuer string
	// AccountName is the user's account identifier (e.g., email).
	AccountName string
	// SecretSize is the size of the secret key in bytes (default: 20).
	SecretSize int
	// Digits is the number of digits in the TOTP code (default: 6).
	Digits int
	// Period is the time period in seconds for TOTP generation (default: 30).
	Period int
	// Algorithm is the hashing algorithm (default: SHA1).
	Algorithm string
}

// DefaultTOTPConfig returns a default TOTP configuration.
func DefaultTOTPConfig() *TOTPConfig {
	return &TOTPConfig{
		Issuer:     "GophKeeper",
		SecretSize: 20,
		Digits:     6,
		Period:     30,
		Algorithm:  "SHA1",
	}
}

// GenerateSecret generates a random secret key for TOTP.
// The secret is returned as a base32-encoded string.
func GenerateSecret(config *TOTPConfig) (string, error) {
	if config == nil {
		config = DefaultTOTPConfig()
	}

	secret := make([]byte, config.SecretSize)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("failed to generate random secret: %w", err)
	}

	return base32.StdEncoding.EncodeToString(secret), nil
}

// GenerateTOTP generates a TOTP code for the given secret.
// The code is valid for the configured time period.
func GenerateTOTP(secret string, config *TOTPConfig) (string, error) {
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP code: %w", err)
	}

	return code, nil
}

// ValidateTOTP validates a TOTP code against the given secret.
// Returns true if the code is valid, false otherwise.
func ValidateTOTP(secret, code string, config *TOTPConfig) bool {
	return totp.Validate(code, secret)
}

// GenerateQRCodeURL generates a URL for QR code generation.
// This URL can be used to create a QR code that users can scan with authenticator apps.
func GenerateQRCodeURL(secret, accountName string, config *TOTPConfig) (string, error) {
	if config == nil {
		config = DefaultTOTPConfig()
	}

	if accountName == "" {
		accountName = config.AccountName
	}

	period := config.Period
	if period < 0 {
		period = 0
	}
	opts := totp.GenerateOpts{
		Issuer:      config.Issuer,
		AccountName: accountName,
		Secret:      []byte(secret),
		Period:      uint(period),
		Digits:      otp.Digits(config.Digits),
		Algorithm:   otp.AlgorithmSHA1,
	}

	key, err := totp.Generate(opts)
	if err != nil {
		return "", err
	}

	return key.URL(), nil
}

// GenerateBackupCodes generates a set of backup codes for account recovery.
// These codes can be used to access the account if the TOTP device is lost.
func GenerateBackupCodes(count int) ([]string, error) {
	if count <= 0 {
		count = 10 // Default to 10 backup codes
	}

	codes := make([]string, count)
	for i := 0; i < count; i++ {
		// Generate 8-digit backup codes
		code := make([]byte, 4) // 4 bytes = 8 hex digits
		if _, err := rand.Read(code); err != nil {
			return nil, fmt.Errorf("failed to generate backup code %d: %w", i+1, err)
		}
		codes[i] = fmt.Sprintf("%08x", code)
	}

	return codes, nil
}

// ValidateBackupCode validates a backup code against a list of valid codes.
// Returns true if the code is valid, false otherwise.
func ValidateBackupCode(code string, validCodes []string) bool {
	for _, validCode := range validCodes {
		if code == validCode {
			return true
		}
	}
	return false
}
