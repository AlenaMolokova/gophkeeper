package api

import (
	"context"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/pkg/errors"
)

// Register handles user registration.
func (s *Server) Register(ctx context.Context, req *clientapi.RegisterRequest) (*clientapi.RegisterResponse, error) {
	token, err := s.auth.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register user")
	}
	return &clientapi.RegisterResponse{Token: token}, nil
}

// Login handles user authentication.
func (s *Server) Login(ctx context.Context, req *clientapi.LoginRequest) (*clientapi.LoginResponse, error) {
	token, err := s.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login user")
	}
	return &clientapi.LoginResponse{Token: token}, nil
}
