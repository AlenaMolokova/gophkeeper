// Package api provides gRPC service implementations for user authentication in the GophKeeper application.
// This file contains the UserServer implementation with methods for user registration
// and login with proper error handling and JWT token generation.
package api

import (
	"context"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/pkg/errors"
)

// Register handles user registration requests.
// It validates the provided email and password, creates a new user account,
// and returns a JWT token for subsequent authenticated requests.
//
// The function uses the underlying auth service to perform the actual registration
// and wraps any errors with additional context for better debugging.
func (s *UserServer) Register(ctx context.Context, req *clientapi.RegisterRequest) (*clientapi.RegisterResponse, error) {
	token, err := s.auth.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register user")
	}
	return &clientapi.RegisterResponse{Token: token}, nil
}

// Login handles user authentication requests.
// It validates the provided email and password against stored credentials,
// and returns a JWT token if authentication is successful.
//
// The function uses the underlying auth service to perform the actual authentication
// and wraps any errors with additional context for better debugging.
func (s *UserServer) Login(ctx context.Context, req *clientapi.LoginRequest) (*clientapi.LoginResponse, error) {
	token, err := s.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login user")
	}
	return &clientapi.LoginResponse{Token: token}, nil
}
