// Package auth provides authentication use cases for the GophKeeper client.
// It implements business logic for user registration and login operations
// following the Single Responsibility Principle and Dependency Inversion.
package auth

import (
	"context"
	"fmt"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// BaseUsecase provides common functionality for authentication use cases.
// It encapsulates shared logic and dependencies for user authentication operations.
type BaseUsecase struct {
	userClient clientapi.UserServiceClient
}

// NewBaseUsecase creates a new BaseUsecase instance.
func NewBaseUsecase(userClient clientapi.UserServiceClient) *BaseUsecase {
	return &BaseUsecase{
		userClient: userClient,
	}
}

// validateCredentials validates that email and password are provided.
// It returns an error if either field is empty.
func (b *BaseUsecase) validateCredentials(email, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("email and password are required")
	}
	return nil
}

// getUserClient returns the user client for use in derived use cases.
func (b *BaseUsecase) getUserClient() clientapi.UserServiceClient {
	return b.userClient
}

// executeLoginOperation executes login operation with common validation logic.
func (b *BaseUsecase) executeLoginOperation(ctx context.Context, email, password string) (string, error) {
	if err := b.validateCredentials(email, password); err != nil {
		return "", err
	}

	resp, err := b.getUserClient().Login(ctx, &clientapi.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	return resp.Token, nil
}

// executeRegisterOperation executes register operation with common validation logic.
func (b *BaseUsecase) executeRegisterOperation(ctx context.Context, email, password string) (string, error) {
	if err := b.validateCredentials(email, password); err != nil {
		return "", err
	}

	resp, err := b.getUserClient().Register(ctx, &clientapi.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("registration failed: %w", err)
	}

	return resp.Token, nil
}
