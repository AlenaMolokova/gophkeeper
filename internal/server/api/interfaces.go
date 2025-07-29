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

	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
)

// DataReader defines the interface for data reading operations.
// This interface follows the Interface Segregation Principle by providing
// only the methods needed for reading data.
type DataReader interface {
	FindDataByID(ctx context.Context, userID, dataID string) (models.Data, error)
}

// DataWriter defines the interface for data writing operations.
// This interface follows the Interface Segregation Principle by providing
// only the methods needed for writing data.
type DataWriter interface {
	SaveData(ctx context.Context, userID string, data models.Data) (string, error)
	EditData(ctx context.Context, userID string, data models.Data) (string, error)
	DeleteData(ctx context.Context, userID, dataID string) error
}

// DataStorage combines DataReader and DataWriter interfaces.
// This interface is used when both reading and writing capabilities are needed.
type DataStorage interface {
	DataReader
	DataWriter
}

//go:generate mockgen -source=interfaces.go -destination=mocks/mocks.go -package=mocks
