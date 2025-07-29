// Package otp provides Time-based One-Time Password (TOTP) functionality for the GophKeeper application.
// It implements TOTP generation, validation, and QR code URL generation for two-factor authentication.
//
// The package supports standard TOTP algorithms and provides backup code generation
// for account recovery scenarios.
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
// It provides customizable settings for TOTP parameters including
// issuer, account name, secret size, digits, period, and algorithm.
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
// It provides sensible defaults for TOTP parameters that are compatible
// with most authenticator applications and follow RFC 6238 standards.
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
// It creates a cryptographically secure random secret of the specified size
// and returns it as a base32-encoded string for compatibility with TOTP standards.
//
// The function uses crypto/rand for secure random number generation
// and ensures the secret is properly encoded for TOTP applications.
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
// It creates a time-based one-time password using the current time
// and the provided secret key, following RFC 6238 standards.
//
// The generated code is valid for the configured time period
// and can be used for two-factor authentication.
func GenerateTOTP(secret string, config *TOTPConfig) (string, error) {
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP code: %w", err)
	}

	return code, nil
}

// ValidateTOTP validates a TOTP code against the given secret.
// It checks if the provided code matches the expected TOTP value
// for the current time window, allowing for clock skew tolerance.
//
// Returns true if the code is valid, false otherwise.
func ValidateTOTP(secret, code string, config *TOTPConfig) bool {
	return totp.Validate(code, secret)
}

// GenerateQRCodeURL generates a URL for QR code generation.
// It creates a standardized URL that can be used to generate QR codes
// for authenticator applications like Google Authenticator or Authy.
//
// The URL follows the otpauth:// scheme and includes all necessary
// parameters for TOTP setup in authenticator apps.
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
	// Safely convert to uint, ensuring it's within valid range
	var periodUint uint
	if period >= 0 && period <= 65535 { // reasonable upper limit for period
		periodUint = uint(period)
	} else {
		periodUint = 30 // default period
	}

	opts := totp.GenerateOpts{
		Issuer:      config.Issuer,
		AccountName: accountName,
		Secret:      []byte(secret),
		Period:      periodUint,
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
// It creates cryptographically secure random codes that can be used
// as an alternative to TOTP codes for account access.
//
// The function generates codes of 8 characters each, using alphanumeric characters
// for easy reading and manual entry if needed.
func GenerateBackupCodes(count int) ([]string, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be positive")
	}

	codes := make([]string, count)
	for i := 0; i < count; i++ {
		code := make([]byte, 8)
		if _, err := rand.Read(code); err != nil {
			return nil, fmt.Errorf("failed to generate backup code: %w", err)
		}
		// Convert to alphanumeric characters
		for j := 0; j < 8; j++ {
			code[j] = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"[code[j]%36]
		}
		codes[i] = string(code)
	}

	return codes, nil
}

// ValidateBackupCode validates a backup code against a list of valid codes.
// It checks if the provided code exists in the list of valid backup codes
// and returns true if the code is valid, false otherwise.
//
// This function is typically used during account recovery or when
// TOTP authentication is not available.
func ValidateBackupCode(code string, validCodes []string) bool {
	for _, validCode := range validCodes {
		if code == validCode {
			return true
		}
	}
	return false
}
