// Package usecases provides business logic implementations for data operations in the GophKeeper application.
// This file contains the AddDataUsecase implementation that handles adding new encrypted data
// with proper authentication and user isolation following Clean Architecture principles.
package usecases

import (
	"context"
	"fmt"

	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
)

// TokenValidator defines the interface for token validation.
type TokenValidator interface {
	ValidateToken(tokenString string) (string, error)
}

// DataSaver defines the interface for saving data.
type DataSaver interface {
	SaveData(ctx context.Context, userID string, data models.Data) (string, error)
}

// AddDataUsecase handles adding new data business logic.
type AddDataUsecase struct {
	tokenValidator TokenValidator
	dataSaver      DataSaver
}

// NewAddDataUsecase creates a new AddDataUsecase instance.
func NewAddDataUsecase(tokenValidator TokenValidator, dataSaver DataSaver) *AddDataUsecase {
	return &AddDataUsecase{
		tokenValidator: tokenValidator,
		dataSaver:      dataSaver,
	}
}

// AddDataRequest represents the add data request.
type AddDataRequest struct {
	Token string
	Data  *models.Data
}

// AddDataResponse represents the add data response.
type AddDataResponse struct {
	ID string
}

// Execute performs adding new data.
func (u *AddDataUsecase) Execute(ctx context.Context, req AddDataRequest) (*AddDataResponse, error) {
	// Validate token and get user ID
	userID, err := u.tokenValidator.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	// Save data
	id, err := u.dataSaver.SaveData(ctx, userID, *req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to save data: %w", err)
	}

	return &AddDataResponse{ID: id}, nil
}
