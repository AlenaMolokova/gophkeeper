// Package auth provides authentication use cases for the GophKeeper client.
// It implements business logic for user registration and login operations
// following the Single Responsibility Principle and Dependency Inversion.
package auth

import (
	"context"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// RegisterUsecase handles user registration business logic.
// It encapsulates the registration flow and provides a clean interface
// for the command-line interface to interact with.
type RegisterUsecase struct {
	*BaseUsecase
}

// NewRegisterUsecase creates a new RegisterUsecase instance.
// It accepts dependencies through the constructor for easy testing
// and dependency injection.
func NewRegisterUsecase(userClient clientapi.UserServiceClient) *RegisterUsecase {
	return &RegisterUsecase{
		BaseUsecase: NewBaseUsecase(userClient),
	}
}

// Execute performs user registration with the provided credentials.
// It validates input parameters and returns the JWT token upon successful registration.
func (r *RegisterUsecase) Execute(ctx context.Context, email, password string) (string, error) {
	return r.executeRegisterOperation(ctx, email, password)
}
