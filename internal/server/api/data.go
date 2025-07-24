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
func (s *Server) validateTokenAndGetUserID(tokenString string) (string, error) {
	userID, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", errors.Wrap(err, "failed to validate token")
	}
	return userID, nil
}

// convertToModelData converts client API data to model data.
func convertToModelData(data *clientapi.Data) models.Data {
	return models.Data{
		ID:        data.Id,
		Type:      data.Type,
		Payload:   data.Payload,
		Metadata:  data.Metadata,
		Timestamp: data.Timestamp,
	}
}

// convertToClientData converts model data to client API data.
func convertToClientData(data models.Data) *clientapi.Data {
	return &clientapi.Data{
		Id:        data.ID,
		Type:      data.Type,
		Payload:   data.Payload,
		Metadata:  data.Metadata,
		Timestamp: data.Timestamp,
	}
}

// AddData handles adding new data.
func (s *Server) AddData(ctx context.Context, req *clientapi.AddDataRequest) (*clientapi.AddDataResponse, error) {
	userID, err := s.validateTokenAndGetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	data := convertToModelData(req.Data)
	id, err := s.storage.SaveData(ctx, userID, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save data")
	}
	return &clientapi.AddDataResponse{Id: id}, nil
}

// GetData handles retrieving data by ID.
func (s *Server) GetData(ctx context.Context, req *clientapi.GetDataRequest) (*clientapi.GetDataResponse, error) {
	userID, err := s.validateTokenAndGetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	data, err := s.storage.FindDataByID(ctx, userID, req.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get data")
	}

	return &clientapi.GetDataResponse{
		Data: convertToClientData(data),
	}, nil
}

// EditData handles updating existing data.
func (s *Server) EditData(ctx context.Context, req *clientapi.EditDataRequest) (*clientapi.EditDataResponse, error) {
	userID, err := s.validateTokenAndGetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	data := convertToModelData(req.Data)
	id, err := s.storage.EditData(ctx, userID, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to edit data")
	}
	return &clientapi.EditDataResponse{Id: id}, nil
}

// DeleteData handles deleting data by ID.
func (s *Server) DeleteData(ctx context.Context, req *clientapi.DeleteDataRequest) (*clientapi.DeleteDataResponse, error) {
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
func (s *Server) ValidateToken(tokenString string) (string, error) {
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
