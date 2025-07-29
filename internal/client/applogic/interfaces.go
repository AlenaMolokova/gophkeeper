// Package applogic provides high-level business logic functions for the GophKeeper client.
// It encapsulates common operations like user authentication, data encryption/decryption,
// and gRPC communication with the server.
//
// The package serves as an abstraction layer between the UI and the underlying
// gRPC clients, providing a clean API for client applications.
package applogic

import (
	"context"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"google.golang.org/grpc"
)

// UserServiceClientInterface defines the interface for user service operations.
// This interface is defined where it's used (next to the consumer) following
// the Interface Segregation Principle.
type UserServiceClientInterface interface {
	Register(ctx context.Context, req *clientapi.RegisterRequest, opts ...grpc.CallOption) (*clientapi.RegisterResponse, error)
	Login(ctx context.Context, req *clientapi.LoginRequest, opts ...grpc.CallOption) (*clientapi.LoginResponse, error)
}

// DataServiceClientInterface defines the interface for data service operations.
// This interface is defined where it's used (next to the consumer) following
// the Interface Segregation Principle.
type DataServiceClientInterface interface {
	AddData(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error)
	GetData(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error)
	EditData(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error)
	DeleteData(ctx context.Context, req *clientapi.DeleteDataRequest, opts ...grpc.CallOption) (*clientapi.DeleteDataResponse, error)
}

//go:generate mockgen -source=interfaces.go -destination=mocks/mocks.go -package=mocks
