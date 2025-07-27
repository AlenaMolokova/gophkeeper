// Package api provides gRPC service implementations for data operations in the GophKeeper application.
// This file contains the DataServer implementation with methods for adding, retrieving,
// editing, and deleting encrypted data with proper authentication and user isolation.
package api

import (
	"context"
	"fmt"

	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// validateTokenAndGetUserID validates the JWT token and returns the user ID.
// This is a helper function that wraps the token validation logic
// and provides consistent error handling across data operations.
func (s *DataServer) validateTokenAndGetUserID(tokenString string) (string, error) {
	userID, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", errors.Wrap(err, "failed to validate token")
	}
	return userID, nil
}

// AddData handles adding new encrypted data for an authenticated user.
// It validates the JWT token, converts the protobuf data to internal model,
// and stores it in the database with proper user isolation.
//
// The function ensures that data is associated with the correct user
// and returns a unique data ID for future reference.
func (s *DataServer) AddData(ctx context.Context, req *clientapi.AddDataRequest) (*clientapi.AddDataResponse, error) {
	userID, err := s.validateTokenAndGetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	data := models.ConvertFromProto(req.Data)
	id, err := s.storage.SaveData(ctx, userID, *data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save data")
	}
	return &clientapi.AddDataResponse{Id: id}, nil
}

// GetData handles retrieving encrypted data by ID for an authenticated user.
// It validates the JWT token and ensures that the requested data belongs
// to the authenticated user, providing proper data isolation.
//
// The function returns the encrypted data with its metadata and timestamp
// for synchronization purposes.
func (s *DataServer) GetData(ctx context.Context, req *clientapi.GetDataRequest) (*clientapi.GetDataResponse, error) {
	userID, err := s.validateTokenAndGetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	data, err := s.storage.FindDataByID(ctx, userID, req.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get data")
	}

	return &clientapi.GetDataResponse{
		Data: data.ConvertToProto(),
	}, nil
}

// EditData handles updating existing encrypted data for an authenticated user.
// It validates the JWT token, ensures the data belongs to the user,
// and updates the data with new content while preserving the data ID.
//
// The function maintains data integrity by verifying ownership
// before allowing any modifications.
func (s *DataServer) EditData(ctx context.Context, req *clientapi.EditDataRequest) (*clientapi.EditDataResponse, error) {
	userID, err := s.validateTokenAndGetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	data := models.ConvertFromProto(req.Data)
	id, err := s.storage.EditData(ctx, userID, *data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to edit data")
	}
	return &clientapi.EditDataResponse{Id: id}, nil
}

// DeleteData handles deleting encrypted data by ID for an authenticated user.
// It validates the JWT token and ensures that the data to be deleted
// belongs to the authenticated user, providing secure data removal.
//
// The function returns an empty response on successful deletion
// or an error if the data doesn't exist or doesn't belong to the user.
func (s *DataServer) DeleteData(ctx context.Context, req *clientapi.DeleteDataRequest) (*clientapi.DeleteDataResponse, error) {
	userID, err := s.validateTokenAndGetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	err = s.storage.DeleteData(ctx, userID, req.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete data")
	}
	return &clientapi.DeleteDataResponse{}, nil
}

// ValidateToken validates the JWT token and returns the user ID.
// It parses the token using the configured JWT secret, verifies the signature,
// and extracts the user ID from the subject claim.
//
// The function ensures that only valid, properly signed tokens
// with the correct algorithm are accepted.
func (s *DataServer) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid user ID in token")
	}
	return userID, nil
}
