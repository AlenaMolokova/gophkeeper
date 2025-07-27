// Package usecases provides business logic implementations for user authentication in the GophKeeper application.
// This file contains the LoginUsecase implementation that handles user login with secure
// password verification and JWT token generation following Clean Architecture principles.
package usecases

import (
	"context"
	"fmt"
)

// UserFinder defines the interface for finding users by email.
// Implementations should return user data including the hashed password
// for subsequent password verification.
type UserFinder interface {
	FindUserByEmail(ctx context.Context, email string) (User, error)
}

// User represents a user entity.
// It contains the user's ID, email, and hashed password for authentication.
type User struct {
	ID       string
	Email    string
	Password string
}

// PasswordVerifier defines the interface for password verification.
// Implementations should compare the provided password against the stored hash
// using secure comparison methods.
type PasswordVerifier interface {
	VerifyPassword(hashedPassword, password string) error
}

// LoginUsecase handles user login business logic.
// It orchestrates the authentication process by finding users,
// verifying passwords, and generating authentication tokens.
type LoginUsecase struct {
	userFinder       UserFinder
	passwordVerifier PasswordVerifier
	tokenGenerator   TokenGenerator
}

// NewLoginUsecase creates a new LoginUsecase instance.
// It accepts dependencies through constructor injection for better testability
// and follows the dependency inversion principle.
func NewLoginUsecase(userFinder UserFinder, passwordVerifier PasswordVerifier, tokenGenerator TokenGenerator) *LoginUsecase {
	return &LoginUsecase{
		userFinder:       userFinder,
		passwordVerifier: passwordVerifier,
		tokenGenerator:   tokenGenerator,
	}
}

// LoginRequest represents the login request.
// It contains the user's email and password for authentication.
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResponse represents the login response.
// It contains the JWT token for authenticated access to the system.
type LoginResponse struct {
	Token string
}

// Execute performs user login.
// It follows a three-step process: find the user by email,
// verify the provided password, and generate an authentication token.
//
// The function ensures secure password verification and returns
// a valid JWT token for authenticated access.
func (u *LoginUsecase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Find user by email
	user, err := u.userFinder.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify password
	err = u.passwordVerifier.VerifyPassword(user.Password, req.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	// Generate JWT token
	token, err := u.tokenGenerator.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{Token: token}, nil
}
