// Package api provides gRPC service implementations for the GophKeeper application.
// It includes UserService and DataService implementations with proper authentication
// and data storage integration.
//
// The package follows Go best practices by defining interfaces where they are used
// and implementing clean separation of concerns between user authentication and
// data management operations.
package api

import (
	"context"

	"github.com/AlenaMolokova/gophkeeper/internal/server/auth"
	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// DataStorage defines the interface for data storage operations.
// This interface is defined where it's used (next to the consumer).
type DataStorage interface {
	SaveData(ctx context.Context, userID string, data models.Data) (string, error)
	FindDataByID(ctx context.Context, userID, dataID string) (models.Data, error)
	EditData(ctx context.Context, userID string, data models.Data) (string, error)
	DeleteData(ctx context.Context, userID, dataID string) error
}

// UserServer implements the UserService gRPC service.
type UserServer struct {
	clientapi.UnimplementedUserServiceServer
	auth      *auth.Auth
	jwtSecret []byte
}

// DataServer implements the DataService gRPC service.
type DataServer struct {
	clientapi.UnimplementedDataServiceServer
	auth      *auth.Auth
	storage   DataStorage
	jwtSecret []byte
}

// NewUserServer creates a new UserServer instance.
func NewUserServer(auth *auth.Auth, jwtSecret []byte) *UserServer {
	return &UserServer{
		auth:      auth,
		jwtSecret: jwtSecret,
	}
}

// NewDataServer creates a new DataServer instance.
func NewDataServer(auth *auth.Auth, storage DataStorage, jwtSecret []byte) *DataServer {
	return &DataServer{
		auth:      auth,
		storage:   storage,
		jwtSecret: jwtSecret,
	}
}
