package sync

import (
	"context"
	"encoding/json"
	"os"

	"github.com/AlenaMolokova/gophkeeper/internal/client/storage"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Syncer defines the interface for data synchronization operations.
// It provides methods for uploading local data to server, downloading from server,
// and merging local and remote data with conflict resolution.
type Syncer interface {
	// Upload uploads local data to the server.
	Upload(ctx context.Context) error
	// Download downloads data from the server to local storage.
	Download(ctx context.Context) error
	// Merge merges local and remote data with conflict resolution.
	Merge(local, remote map[string][]byte) map[string][]byte
}

// DefaultSyncer implements the Syncer interface with default storage.
type DefaultSyncer struct {
	Storage *storage.Storage
}

// NewDefaultSyncer creates a new syncer with default storage path.
func NewDefaultSyncer() (*DefaultSyncer, error) {
	dbPath := storage.GetDefaultDBPath()
	st, err := storage.NewStorage(dbPath)
	if err != nil {
		return nil, err
	}
	return &DefaultSyncer{Storage: st}, nil
}

// NewDefaultSyncerWithPath creates a new syncer with custom storage path.
func NewDefaultSyncerWithPath(dbPath string) (*DefaultSyncer, error) {
	st, err := storage.NewStorage(dbPath)
	if err != nil {
		return nil, err
	}
	return &DefaultSyncer{Storage: st}, nil
}

func getToken() string {
	token := os.Getenv("GOPHKEEPER_TOKEN")
	if token == "" {
		panic("GOPHKEEPER_TOKEN не установлен")
	}
	return token
}

func getCertPath() string {
	certPath := os.Getenv("GOPHKEEPER_CERT")
	if certPath == "" {
		certPath = "cert/server.crt"
	}
	return certPath
}

func (s *DefaultSyncer) Upload(ctx context.Context) error {
	creds, err := credentials.NewClientTLSFromFile(getCertPath(), "")
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := clientapi.NewGophKeeperClient(conn)

	all, err := s.Storage.GetAllData()
	if err != nil {
		return err
	}
	token := getToken()
	for _, raw := range all {
		var data clientapi.Data
		if err := json.Unmarshal(raw, &data); err != nil {
			continue
		}
		// Try to update, if doesn't exist — add
		_, err := client.EditData(ctx, &clientapi.EditDataRequest{Token: token, Data: &data})
		if err != nil {
			_, err2 := client.AddData(ctx, &clientapi.AddDataRequest{Token: token, Data: &data})
			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}

func (s *DefaultSyncer) Download(ctx context.Context) error {
	creds, err := credentials.NewClientTLSFromFile(getCertPath(), "")
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := clientapi.NewGophKeeperClient(conn)

	// For simplicity: list of ids should be known locally (or store map id->timestamp)
	local, err := s.Storage.GetAllData()
	if err != nil {
		return err
	}
	token := getToken()
	for id := range local {
		resp, err := client.GetData(ctx, &clientapi.GetDataRequest{Token: token, Id: id})
		if err != nil {
			continue
		}
		remoteData := resp.Data
		var localData clientapi.Data
		_ = json.Unmarshal(local[id], &localData)
		if remoteData.Timestamp > localData.Timestamp {
			b, _ := json.Marshal(remoteData)
			if err := s.Storage.SaveData(id, b); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *DefaultSyncer) Merge(local, remote map[string][]byte) map[string][]byte {
	result := make(map[string][]byte)
	for id, lraw := range local {
		result[id] = lraw
	}
	for id, rraw := range remote {
		if lraw, ok := result[id]; ok {
			var ldata, rdata clientapi.Data
			_ = json.Unmarshal(lraw, &ldata)
			_ = json.Unmarshal(rraw, &rdata)
			if rdata.Timestamp > ldata.Timestamp {
				result[id] = rraw
			}
		} else {
			result[id] = rraw
		}
	}
	return result
}
