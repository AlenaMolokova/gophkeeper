// Package adapters provides adapter implementations for the GophKeeper application.
// This file contains adapters that bridge between usecase interfaces and concrete
// implementations, following the Adapter pattern for dependency inversion.
package adapters

import (
	"context"
	"time"

	"github.com/AlenaMolokova/gophkeeper/internal/server/api/usecases"
	"github.com/AlenaMolokova/gophkeeper/internal/server/auth"
	"github.com/AlenaMolokova/gophkeeper/internal/server/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthAdapter adapts the existing auth.Auth to usecase interfaces.
type AuthAdapter struct {
	auth    *auth.Auth
	storage *storage.PostgresStorage
}

// NewAuthAdapter creates a new AuthAdapter instance.
func NewAuthAdapter(auth *auth.Auth, storage *storage.PostgresStorage) *AuthAdapter {
	return &AuthAdapter{auth: auth, storage: storage}
}

// UserSaverAdapter adapts storage.Storage to UserSaver interface.
type UserSaverAdapter struct {
	storage *storage.PostgresStorage
}

// NewUserSaverAdapter creates a new UserSaverAdapter instance.
func NewUserSaverAdapter(storage *storage.PostgresStorage) *UserSaverAdapter {
	return &UserSaverAdapter{storage: storage}
}

// SaveUser implements usecases.UserSaver interface.
func (a *UserSaverAdapter) SaveUser(ctx context.Context, email, hash string) (string, error) {
	return a.storage.SaveUser(ctx, email, hash)
}

// UserFinderAdapter adapts storage.Storage to UserFinder interface.
type UserFinderAdapter struct {
	storage *storage.PostgresStorage
}

// NewUserFinderAdapter creates a new UserFinderAdapter instance.
func NewUserFinderAdapter(storage *storage.PostgresStorage) *UserFinderAdapter {
	return &UserFinderAdapter{storage: storage}
}

// FindUserByEmail implements usecases.UserFinder interface.
func (a *UserFinderAdapter) FindUserByEmail(ctx context.Context, email string) (usecases.User, error) {
	user, err := a.storage.FindUserByEmail(ctx, email)
	if err != nil {
		return usecases.User{}, err
	}
	return usecases.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Hash, // Note: using Hash field from models.User
	}, nil
}

// PasswordHasherAdapter provides password hashing functionality.
type PasswordHasherAdapter struct{}

// NewPasswordHasherAdapter creates a new PasswordHasherAdapter instance.
func NewPasswordHasherAdapter() *PasswordHasherAdapter {
	return &PasswordHasherAdapter{}
}

// HashPassword implements usecases.PasswordHasher interface.
func (a *PasswordHasherAdapter) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}
	return string(hash), nil
}

// PasswordVerifierAdapter provides password verification functionality.
type PasswordVerifierAdapter struct{}

// NewPasswordVerifierAdapter creates a new PasswordVerifierAdapter instance.
func NewPasswordVerifierAdapter() *PasswordVerifierAdapter {
	return &PasswordVerifierAdapter{}
}

// VerifyPassword implements usecases.PasswordVerifier interface.
func (a *PasswordVerifierAdapter) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// TokenGeneratorAdapter provides JWT token generation functionality.
type TokenGeneratorAdapter struct {
	jwtSecret string
	tokenTTL  time.Duration
}

// NewTokenGeneratorAdapter creates a new TokenGeneratorAdapter instance.
func NewTokenGeneratorAdapter(jwtSecret string, tokenTTL time.Duration) *TokenGeneratorAdapter {
	return &TokenGeneratorAdapter{
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

// GenerateToken implements usecases.TokenGenerator interface.
func (a *TokenGeneratorAdapter) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(a.tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}
