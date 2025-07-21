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

func (s *Storage) SaveData(id string, data []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("data"))
		return b.Put([]byte(id), data)
	})
}

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

func (s *Storage) DeleteData(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("data"))
		return b.Delete([]byte(id))
	})
}

func (s *Storage) Close() error {
	return s.db.Close()
}
