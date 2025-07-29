// Package api provides gRPC service implementations for the GophKeeper application.
// It includes UserService and DataService implementations with proper authentication
// and data storage integration.
//
// The package follows Go best practices by defining interfaces where they are used
// and implementing clean separation of concerns between user authentication and
// data management operations.
package api

import (
	"github.com/AlenaMolokova/gophkeeper/internal/server/auth"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

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
	reader    DataReader
	writer    DataWriter
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
		reader:    storage,
		writer:    storage,
		jwtSecret: jwtSecret,
	}
}
