// Package usecases provides business logic implementations for data operations in the GophKeeper application.
// This file contains the DeleteDataUsecase implementation that handles removing encrypted data
// with proper authentication and user isolation following Clean Architecture principles.
package usecases

import (
	"context"
	"fmt"
)

// DataDeleter defines the interface for deleting data.
type DataDeleter interface {
	DeleteData(ctx context.Context, userID, dataID string) error
}

// DeleteDataUsecase handles deleting data business logic.
type DeleteDataUsecase struct {
	tokenValidator TokenValidator
	dataDeleter    DataDeleter
}

// NewDeleteDataUsecase creates a new DeleteDataUsecase instance.
func NewDeleteDataUsecase(tokenValidator TokenValidator, dataDeleter DataDeleter) *DeleteDataUsecase {
	return &DeleteDataUsecase{
		tokenValidator: tokenValidator,
		dataDeleter:    dataDeleter,
	}
}

// DeleteDataRequest represents the delete data request.
type DeleteDataRequest struct {
	Token string
	ID    string
}

// DeleteDataResponse represents the delete data response.
type DeleteDataResponse struct{}

// Execute performs deleting data.
func (u *DeleteDataUsecase) Execute(ctx context.Context, req DeleteDataRequest) (*DeleteDataResponse, error) {
	// Validate token and get user ID
	userID, err := u.tokenValidator.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	// Delete data
	err = u.dataDeleter.DeleteData(ctx, userID, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete data: %w", err)
	}

	return &DeleteDataResponse{}, nil
}
