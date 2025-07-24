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
type Storage interface {
	SaveUser(ctx context.Context, email, hash string) (string, error)
	FindUserByEmail(ctx context.Context, email string) (models.User, error)
}

// Auth manages user authentication.
type Auth struct {
	storage   Storage
	jwtSecret string
	tokenTTL  time.Duration
}

// NewAuth creates a new Auth instance.
func NewAuth(storage Storage, jwtSecret string, tokenTTL time.Duration) *Auth {
	return &Auth{
		storage:   storage,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

// Register registers a new user and returns a JWT token.
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
func (a *Auth) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(a.tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}
