// Package sync provides data synchronization functionality for the GophKeeper client.
// It handles bidirectional synchronization between local storage and the remote server,
// including conflict resolution and data merging capabilities.
//
// The package implements a simple but effective synchronization strategy
// that ensures data consistency across multiple client instances.
package sync

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AlenaMolokova/gophkeeper/internal/client/storage"
	"github.com/AlenaMolokova/gophkeeper/internal/config"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Syncer defines the interface for data synchronization.
// It provides methods for uploading local data to the server,
// downloading data from the server, and merging conflicting data.
type Syncer interface {
	// Upload uploads local data to the server.
	// It attempts to update existing data or create new entries
	// if the data doesn't exist on the server.
	Upload(ctx context.Context) error
	// Download downloads data from the server to local storage.
	// It retrieves all data associated with the authenticated user
	// and stores it locally for offline access.
	Download(ctx context.Context) error
	// Merge merges local and remote data with conflict resolution.
	// It takes local and remote data maps and returns a merged result
	// based on the implemented conflict resolution strategy.
	Merge(local, remote map[string][]byte) map[string][]byte
}

// DefaultSyncer implements the Syncer interface with default storage.
// It provides a complete synchronization implementation using local SQLite storage
// and gRPC communication with the remote server.
type DefaultSyncer struct {
	Storage *storage.Storage     // Local storage for data persistence
	config  *config.ClientConfig // Client configuration for server connection
}

// NewDefaultSyncer creates a new DefaultSyncer with default storage path.
// It initializes the syncer with the default SQLite database location
// and default client configuration.
func NewDefaultSyncer() (*DefaultSyncer, error) {
	storage, err := storage.NewStorage(storage.GetDefaultDBPath())
	if err != nil {
		return nil, err
	}
	return &DefaultSyncer{
		Storage: storage,
		config:  config.NewClientConfig(),
	}, nil
}

// NewDefaultSyncerWithPath creates a new DefaultSyncer with custom storage path.
// It allows specifying a custom database location for the local storage,
// useful for testing or when the default location is not suitable.
func NewDefaultSyncerWithPath(dbPath string) (*DefaultSyncer, error) {
	storage, err := storage.NewStorage(dbPath)
	if err != nil {
		return nil, err
	}
	return &DefaultSyncer{
		Storage: storage,
		config:  config.NewClientConfig(),
	}, nil
}

// NewDefaultSyncerWithConfig creates a new DefaultSyncer with custom configuration.
// It allows full customization of both storage path and client configuration,
// providing maximum flexibility for different deployment scenarios.
func NewDefaultSyncerWithConfig(dbPath string, cfg *config.ClientConfig) (*DefaultSyncer, error) {
	storage, err := storage.NewStorage(dbPath)
	if err != nil {
		return nil, err
	}
	return &DefaultSyncer{
		Storage: storage,
		config:  cfg,
	}, nil
}

// Upload uploads local data to the server.
// It establishes a secure connection to the server, retrieves all local data,
// and attempts to update existing entries or create new ones as needed.
//
// The function uses a "try update, then add" strategy to handle both
// existing and new data entries efficiently.
func (s *DefaultSyncer) Upload(ctx context.Context) error {
	creds, err := credentials.NewClientTLSFromFile(s.config.CertPath, "")
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(s.config.GetClientServerAddress(), grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	dataClient := clientapi.NewDataServiceClient(conn)

	all, err := s.Storage.GetAllData()
	if err != nil {
		return err
	}
	token := s.config.Token
	for _, raw := range all {
		var data clientapi.Data
		if err := json.Unmarshal(raw, &data); err != nil {
			continue
		}
		// Try to update, if doesn't exist — add
		_, err := dataClient.EditData(ctx, &clientapi.EditDataRequest{Token: token, Data: &data})
		if err != nil {
			_, err2 := dataClient.AddData(ctx, &clientapi.AddDataRequest{Token: token, Data: &data})
			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}

// Download downloads data from the server to local storage.
// It establishes a secure connection to the server and retrieves
// all data associated with the authenticated user for local storage.
//
// The function downloads data for all locally known IDs and updates
// the local storage with the latest server data.
func (s *DefaultSyncer) Download(ctx context.Context) error {
	creds, err := credentials.NewClientTLSFromFile(s.config.CertPath, "")
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(s.config.GetClientServerAddress(), grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	dataClient := clientapi.NewDataServiceClient(conn)

	// For simplicity: list of ids should be known locally (or store map id->timestamp)
	local, err := s.Storage.GetAllData()
	if err != nil {
		return err
	}
	token := s.config.Token
	for id := range local {
		resp, err := dataClient.GetData(ctx, &clientapi.GetDataRequest{Token: token, Id: id})
		if err != nil {
			continue
		}
		// Update local data
		raw, err := json.Marshal(resp.Data)
		if err != nil {
			continue
		}
		err = s.Storage.SaveData(id, raw)
		if err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}
	}
	return nil
}

// Merge merges local and remote data with conflict resolution.
// It implements a simple merge strategy where remote data takes precedence
// over local data in case of conflicts.
//
// The function creates a new map containing all data from both sources,
// with remote data overwriting local data for duplicate keys.
func (s *DefaultSyncer) Merge(local, remote map[string][]byte) map[string][]byte {
	// Simple merge strategy: remote wins
	result := make(map[string][]byte)
	for k, v := range local {
		result[k] = v
	}
	for k, v := range remote {
		result[k] = v
	}
	return result
}
