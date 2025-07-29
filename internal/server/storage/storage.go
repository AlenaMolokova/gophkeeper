// Package storage provides database storage functionality for the GophKeeper server.
// It includes PostgreSQL connection pooling, user management, and data storage operations.
//
// The package supports configurable connection pools with monitoring capabilities
// and provides both synchronous and asynchronous database operations.
//
// Example usage:
//
//	config := storage.DefaultPoolConfig()
//	config.MaxConns = 50
//	storage, err := storage.NewPostgresStorageWithConfig(ctx, dsn, config)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer storage.Close(ctx)
//
//go:generate mockgen -source=storage.go -destination=storage_mock.go -package=storage

package storage

import (
	"context"
	"time"

	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// PoolConfig holds configuration for the database connection pool.
// It allows fine-tuning of connection pool behavior for optimal performance
// in different deployment scenarios (development, production, high-load).
type PoolConfig struct {
	// MaxConns is the maximum number of connections in the pool.
	// Should not exceed PostgreSQL max_connections setting.
	MaxConns int32

	// MinConns is the minimum number of connections to maintain in the pool.
	// These connections are created at startup and kept alive.
	MinConns int32

	// MaxConnLifetime is the maximum time a connection can be reused.
	// After this time, connections are closed and recreated.
	MaxConnLifetime time.Duration

	// MaxConnIdleTime is the maximum time a connection can remain idle.
	// Idle connections beyond this limit are closed.
	MaxConnIdleTime time.Duration

	// HealthCheckPeriod is the interval between health checks of connections.
	// Shorter intervals provide faster detection of connection issues.
	HealthCheckPeriod time.Duration
}

// DefaultPoolConfig returns default pool configuration suitable for most use cases.
// The default values are optimized for moderate load scenarios and can be
// customized through environment variables or direct configuration.
//
// Default values:
//   - MaxConns: 20
//   - MinConns: 5
//   - MaxConnLifetime: 1 hour
//   - MaxConnIdleTime: 30 minutes
//   - HealthCheckPeriod: 1 minute
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		MaxConns:          20,
		MinConns:          5,
		MaxConnLifetime:   1 * time.Hour,
		MaxConnIdleTime:   30 * time.Minute,
		HealthCheckPeriod: 1 * time.Minute,
	}
}

// PoolStats holds statistics about the connection pool.
// These metrics help monitor pool health and performance.
type PoolStats struct {
	// TotalConns is the total number of connections in the pool.
	TotalConns int32

	// IdleConns is the number of connections currently available for use.
	IdleConns int32

	// AcquiredConns is the number of connections currently in use.
	AcquiredConns int32

	// Constructing is the number of connections being created.
	Constructing int32
}

// PostgresStorage implements the Storage interface using PostgreSQL.
// It manages a connection pool for efficient database operations and provides
// thread-safe access to PostgreSQL database.
type PostgresStorage struct {
	// pool is the underlying PostgreSQL connection pool.
	pool *pgxpool.Pool
}

// NewPostgresStorage creates a new PostgresStorage instance with default pool configuration.
// This is the recommended way to create a storage instance for most use cases.
//
// The function uses DefaultPoolConfig() for connection pool settings,
// which provides reasonable defaults for moderate load scenarios.
func NewPostgresStorage(ctx context.Context, dsn string) (*PostgresStorage, error) {
	return NewPostgresStorageWithConfig(ctx, dsn, DefaultPoolConfig())
}

// NewPostgresStorageWithConfig creates a new PostgresStorage instance with custom pool configuration.
// Use this function when you need fine-grained control over connection pool behavior.
//
// The function validates the connection string, applies the provided configuration,
// creates the connection pool, and verifies connectivity before returning.
//
// Example:
//
//	config := &storage.PoolConfig{
//		MaxConns: 50,
//		MinConns: 10,
//		MaxConnLifetime: 2 * time.Hour,
//	}
//	storage, err := storage.NewPostgresStorageWithConfig(ctx, dsn, config)
func NewPostgresStorageWithConfig(ctx context.Context, dsn string, config *PoolConfig) (*PostgresStorage, error) {
	if config == nil {
		return nil, errors.New("pool configuration cannot be nil")
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse connection string")
	}

	// Validate pool configuration
	if config.MaxConns <= 0 {
		return nil, errors.New("MaxConns must be positive")
	}
	if config.MinConns < 0 {
		return nil, errors.New("MinConns cannot be negative")
	}
	if config.MinConns > config.MaxConns {
		return nil, errors.New("MinConns cannot be greater than MaxConns")
	}
	if config.MaxConnLifetime < 0 {
		return nil, errors.New("MaxConnLifetime cannot be negative")
	}
	if config.MaxConnIdleTime < 0 {
		return nil, errors.New("MaxConnIdleTime cannot be negative")
	}
	if config.HealthCheckPeriod < 0 {
		return nil, errors.New("HealthCheckPeriod cannot be negative")
	}

	// Apply custom pool configuration
	poolConfig.MaxConns = config.MaxConns
	poolConfig.MinConns = config.MinConns
	poolConfig.MaxConnLifetime = config.MaxConnLifetime
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = config.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create connection pool")
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return &PostgresStorage{pool: pool}, nil
}

// SaveUser creates a new user in the database and returns the user ID.
// The email must be unique across all users.
// If the email already exists, returns an error.
func (s *PostgresStorage) SaveUser(ctx context.Context, email, hash string) (string, error) {
	var id string
	err := s.pool.QueryRow(ctx, SaveUserQuery, email, hash).Scan(&id)
	if err != nil {
		return "", errors.Wrap(err, "failed to save user")
	}
	return id, nil
}

// FindUserByEmail retrieves a user by email.
// Returns a models.User struct containing the user's ID, email, and password hash.
// If no user is found with the given email, returns an error.
func (s *PostgresStorage) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := s.pool.QueryRow(ctx, FindUserByEmailQuery, email).Scan(&user.ID, &user.Email, &user.Hash)
	if err != nil {
		return models.User{}, errors.Wrap(err, "failed to find user")
	}
	return user, nil
}

// SaveData saves a data entry for a user.
// The data is associated with the specified userID for proper isolation.
// The payload should be encrypted before calling this function.
// Returns the unique data ID for future reference.
func (s *PostgresStorage) SaveData(ctx context.Context, userID string, data models.Data) (string, error) {
	var id string
	err := s.pool.QueryRow(ctx, SaveDataQuery, userID, data.Type, data.Payload, data.Metadata, data.Timestamp).Scan(&id)
	if err != nil {
		return "", errors.Wrap(err, "failed to save data")
	}
	return id, nil
}

// FindDataByID retrieves a data entry by ID for a user.
// Ensures data isolation by verifying the data belongs to the specified user.
// Returns the encrypted data if found, or an error if not found or access denied.
func (s *PostgresStorage) FindDataByID(ctx context.Context, userID, dataID string) (models.Data, error) {
	var data models.Data
	err := s.pool.QueryRow(ctx, FindDataByIDQuery, dataID, userID).
		Scan(&data.ID, &data.Type, &data.Payload, &data.Metadata, &data.Timestamp)
	if err != nil {
		return models.Data{}, errors.Wrap(err, "failed to find data")
	}
	return data, nil
}

// EditData updates a data entry for a user.
// The data ID must exist and belong to the specified user.
// The payload should be encrypted before calling this function.
// Returns the data ID upon successful update.
func (s *PostgresStorage) EditData(ctx context.Context, userID string, data models.Data) (string, error) {
	var id string
	err := s.pool.QueryRow(ctx, EditDataQuery, data.ID, userID, data.Type, data.Payload, data.Metadata, data.Timestamp).Scan(&id)
	if err != nil {
		return "", errors.Wrap(err, "failed to edit data")
	}
	return id, nil
}

// DeleteData deletes a data entry for a user.
// Ensures data isolation by verifying the data belongs to the specified user.
// Returns an error if the data doesn't exist or doesn't belong to the user.
func (s *PostgresStorage) DeleteData(ctx context.Context, userID, dataID string) error {
	result, err := s.pool.Exec(ctx, DeleteDataQuery, dataID, userID)
	if err != nil {
		return errors.Wrap(err, "failed to delete data")
	}
	if result.RowsAffected() == 0 {
		return errors.New("no data found for deletion")
	}
	return nil
}

// GetPoolStats returns current pool statistics.
// This method provides real-time metrics about the connection pool's state,
// useful for monitoring, debugging, and performance optimization.
//
// The returned statistics include:
// - Total connections in the pool.
// - Idle connections available for use.
// - Acquired connections currently in use.
// - Connections being constructed.
func (s *PostgresStorage) GetPoolStats() *PoolStats {
	stats := s.pool.Stat()
	return &PoolStats{
		TotalConns:    stats.TotalConns(),
		IdleConns:     stats.IdleConns(),
		AcquiredConns: stats.AcquiredConns(),
		Constructing:  stats.ConstructingConns(),
	}
}

// Close closes the database connection pool.
// This method should be called when shutting down the application
// to properly release all database connections and resources.
// It's safe to call this method multiple times.
func (s *PostgresStorage) Close(ctx context.Context) error {
	s.pool.Close()
	return nil
}
