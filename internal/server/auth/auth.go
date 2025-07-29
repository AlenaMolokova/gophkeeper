// Package auth provides user authentication functionality for the GophKeeper server.
// It handles user registration, login, password hashing, and JWT token generation
// with secure bcrypt password hashing and configurable token expiration.
//
// The package implements secure authentication practices including
// password hashing with bcrypt and JWT token generation with expiration.
package auth

import (
	"context"
	"time"

	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Storage defines the interface for auth storage operations.
// It provides methods for saving and finding users by email,
// allowing the auth package to work with any storage implementation.
type Storage interface {
	SaveUser(ctx context.Context, email, hash string) (string, error)
	FindUserByEmail(ctx context.Context, email string) (models.User, error)
}

// Auth manages user authentication.
// It provides user registration, login, and token generation functionality
// with secure password handling and JWT token management.
type Auth struct {
	storage   Storage       // Storage interface for user persistence
	jwtSecret string        // Secret key for JWT token signing
	tokenTTL  time.Duration // Token time-to-live duration
}

// NewAuth creates a new Auth instance.
// It initializes the authentication service with the provided storage,
// JWT secret, and token expiration time.
func NewAuth(storage Storage, jwtSecret string, tokenTTL time.Duration) *Auth {
	return &Auth{
		storage:   storage,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

// Register registers a new user and returns a JWT token.
// It hashes the provided password using bcrypt, saves the user to storage,
// and generates a JWT token for immediate authentication.
//
// The function ensures that passwords are securely hashed before storage
// and returns a valid JWT token upon successful registration.
func (a *Auth) Register(ctx context.Context, email, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}

	userID, err := a.storage.SaveUser(ctx, email, string(hash))
	if err != nil {
		return "", errors.Wrap(err, "failed to save user")
	}

	token, err := a.generateToken(userID)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}
	return token, nil
}

// Login authenticates a user and returns a JWT token.
// It retrieves the user from storage, verifies the password using bcrypt,
// and generates a JWT token upon successful authentication.
//
// The function ensures secure password verification and returns
// a valid JWT token for authenticated access.
func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.storage.FindUserByEmail(ctx, email)
	if err != nil {
		return "", errors.Wrap(err, "failed to find user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil {
		return "", errors.Wrap(err, "invalid password")
	}

	token, err := a.generateToken(user.ID)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}
	return token, nil
}

// generateToken generates a JWT token for a user.
// It creates a JWT token with the user ID as the subject claim
// and an expiration time based on the configured token TTL.
//
// The function uses HMAC-SHA256 signing method and includes
// standard JWT claims for security and validation.
func (a *Auth) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(a.tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}
