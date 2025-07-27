// Package usecases provides business logic implementations for the GophKeeper application.
// It contains usecase structures that encapsulate specific business operations
// and follow the Clean Architecture principles with dependency injection.
//
// Each usecase is responsible for a single business operation and uses
// interfaces defined at the point of use for maximum flexibility and testability.
package usecases

import (
	"context"
	"fmt"
)

// UserSaver defines the interface for saving users.
// This interface is defined where it's used (next to the consumer).
type UserSaver interface {
	SaveUser(ctx context.Context, email, hash string) (string, error)
}

// PasswordHasher defines the interface for password hashing.
// Implementations should use secure hashing algorithms like bcrypt.
type PasswordHasher interface {
	HashPassword(password string) (string, error)
}

// TokenGenerator defines the interface for generating JWT tokens.
// Implementations should create properly signed tokens with appropriate claims.
type TokenGenerator interface {
	GenerateToken(userID string) (string, error)
}

// RegisterUsecase handles user registration business logic.
// It orchestrates the registration process by hashing passwords,
// saving user data, and generating authentication tokens.
type RegisterUsecase struct {
	userSaver      UserSaver
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
}

// NewRegisterUsecase creates a new RegisterUsecase instance.
// It accepts dependencies through constructor injection for better testability
// and follows the dependency inversion principle.
func NewRegisterUsecase(userSaver UserSaver, passwordHasher PasswordHasher, tokenGenerator TokenGenerator) *RegisterUsecase {
	return &RegisterUsecase{
		userSaver:      userSaver,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

// RegisterRequest represents the registration request.
// It contains the user's email and password for account creation.
type RegisterRequest struct {
	Email    string
	Password string
}

// RegisterResponse represents the registration response.
// It contains the JWT token for immediate authentication after registration.
type RegisterResponse struct {
	Token string
}

// Execute performs user registration.
// It follows a three-step process: hash the password, save the user,
// and generate an authentication token.
//
// The function ensures that passwords are properly hashed before storage
// and returns a valid JWT token for immediate use.
func (u *RegisterUsecase) Execute(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// Hash the password
	hashedPassword, err := u.passwordHasher.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Save the user
	userID, err := u.userSaver.SaveUser(ctx, req.Email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Generate JWT token
	token, err := u.tokenGenerator.GenerateToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &RegisterResponse{Token: token}, nil
}
