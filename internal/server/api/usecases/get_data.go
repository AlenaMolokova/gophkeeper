// Package usecases provides business logic implementations for data operations in the GophKeeper application.
// This file contains the GetDataUsecase implementation that handles retrieving encrypted data
// with proper authentication and user isolation following Clean Architecture principles.
package usecases

import (
	"context"
	"fmt"

	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
)

// DataFinder defines the interface for finding data by ID.
type DataFinder interface {
	FindDataByID(ctx context.Context, userID, dataID string) (models.Data, error)
}

// GetDataUsecase handles getting data business logic.
type GetDataUsecase struct {
	tokenValidator TokenValidator
	dataFinder     DataFinder
}

// NewGetDataUsecase creates a new GetDataUsecase instance.
func NewGetDataUsecase(tokenValidator TokenValidator, dataFinder DataFinder) *GetDataUsecase {
	return &GetDataUsecase{
		tokenValidator: tokenValidator,
		dataFinder:     dataFinder,
	}
}

// GetDataRequest represents the get data request.
type GetDataRequest struct {
	Token string
	ID    string
}

// GetDataResponse represents the get data response.
type GetDataResponse struct {
	Data *models.Data
}

// Execute performs getting data by ID.
func (u *GetDataUsecase) Execute(ctx context.Context, req GetDataRequest) (*GetDataResponse, error) {
	// Validate token and get user ID
	userID, err := u.tokenValidator.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	// Find data by ID
	data, err := u.dataFinder.FindDataByID(ctx, userID, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	return &GetDataResponse{Data: &data}, nil
}
