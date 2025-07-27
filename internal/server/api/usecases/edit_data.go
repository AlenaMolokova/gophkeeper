// Package usecases provides business logic implementations for data operations in the GophKeeper application.
// This file contains the EditDataUsecase implementation that handles updating encrypted data
// with proper authentication and user isolation following Clean Architecture principles.
package usecases

import (
	"context"
	"fmt"

	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
)

// DataEditor defines the interface for editing data.
type DataEditor interface {
	EditData(ctx context.Context, userID string, data models.Data) (string, error)
}

// EditDataUsecase handles editing data business logic.
type EditDataUsecase struct {
	tokenValidator TokenValidator
	dataEditor     DataEditor
}

// NewEditDataUsecase creates a new EditDataUsecase instance.
func NewEditDataUsecase(tokenValidator TokenValidator, dataEditor DataEditor) *EditDataUsecase {
	return &EditDataUsecase{
		tokenValidator: tokenValidator,
		dataEditor:     dataEditor,
	}
}

// EditDataRequest represents the edit data request.
type EditDataRequest struct {
	Token string
	Data  *models.Data
}

// EditDataResponse represents the edit data response.
type EditDataResponse struct {
	ID string
}

// Execute performs editing data.
func (u *EditDataUsecase) Execute(ctx context.Context, req EditDataRequest) (*EditDataResponse, error) {
	// Validate token and get user ID
	userID, err := u.tokenValidator.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	// Edit data
	id, err := u.dataEditor.EditData(ctx, userID, *req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to edit data: %w", err)
	}

	return &EditDataResponse{ID: id}, nil
}
