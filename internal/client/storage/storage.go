package storage

import (
	"errors"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
)

// GetDefaultDBPath returns the default path for the bbolt database.
func GetDefaultDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory cannot be determined
		return "gophkeeper.db"
	}
	return filepath.Join(homeDir, ".gophkeeper", "data.db")
}

// EnsureDBDirectory creates the directory for the database if it doesn't exist.
func EnsureDBDirectory(dbPath string) error {
	dir := filepath.Dir(dbPath)
	return os.MkdirAll(dir, 0o700)
}

// Storage provides local data storage using bbolt database.
// It implements a simple key-value storage system for encrypted data
// with automatic directory creation and proper file permissions.
type Storage struct {
	db *bbolt.DB
}

// NewStorage creates a new storage instance with the given database path.
func NewStorage(path string) (*Storage, error) {
	// Ensure the directory exists
	if err := EnsureDBDirectory(path); err != nil {
		return nil, err
	}

	db, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, err
	}

	// Create bucket for data if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("data"))
		return err
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Storage{db: db}, nil
}

// SaveData stores encrypted data with the given ID.
// It saves the data to the bbolt database in a transaction-safe manner.
func (s *Storage) SaveData(id string, data []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("data"))
		return b.Put([]byte(id), data)
	})
}

// GetData retrieves encrypted data by ID.
// It returns the data as a byte slice or an error if the data is not found.
func (s *Storage) GetData(id string) ([]byte, error) {
	var data []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("data"))
		v := b.Get([]byte(id))
		if v == nil {
			return errors.New("not found")
		}
		data = append([]byte{}, v...)
		return nil
	})
	return data, err
}

// GetAllData retrieves all stored data as a map of ID to data.
// It returns all key-value pairs stored in the database for synchronization purposes.
func (s *Storage) GetAllData() (map[string][]byte, error) {
	result := make(map[string][]byte)
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("data"))
		return b.ForEach(func(k, v []byte) error {
			result[string(k)] = append([]byte{}, v...)
			return nil
		})
	})
	return result, err
}

// DeleteData removes data with the given ID from storage.
// It deletes the specified data from the bbolt database in a transaction-safe manner.
func (s *Storage) DeleteData(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("data"))
		return b.Delete([]byte(id))
	})
}

// Close closes the database connection and releases associated resources.
// It should be called when the storage is no longer needed to prevent resource leaks.
func (s *Storage) Close() error {
	return s.db.Close()
}
