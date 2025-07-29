// Package auth provides authentication use cases for the GophKeeper client.
// It implements business logic for user registration and login operations
// following the Single Responsibility Principle and Dependency Inversion.
package auth

import (
	"context"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// LoginUsecase handles user login business logic.
// It encapsulates the login flow and provides a clean interface
// for the command-line interface to interact with.
type LoginUsecase struct {
	*BaseUsecase
}

// NewLoginUsecase creates a new LoginUsecase instance.
// It accepts dependencies through the constructor for easy testing
// and dependency injection.
func NewLoginUsecase(userClient clientapi.UserServiceClient) *LoginUsecase {
	return &LoginUsecase{
		BaseUsecase: NewBaseUsecase(userClient),
	}
}

// Execute performs user login with the provided credentials.
// It validates input parameters and returns the JWT token upon successful authentication.
func (l *LoginUsecase) Execute(ctx context.Context, email, password string) (string, error) {
	return l.executeLoginOperation(ctx, email, password)
}
