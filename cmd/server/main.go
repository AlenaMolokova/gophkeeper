// Package main provides the GophKeeper server application.
// The server implements a secure password manager with gRPC API,
// PostgreSQL storage, and configurable connection pooling.
//
// The server supports:
// - JWT-based authentication
// - Encrypted data storage
// - Configurable database connection pools
// - TLS-secured gRPC communication
// - Real-time pool monitoring
package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/AlenaMolokova/gophkeeper/internal/config"
	"github.com/AlenaMolokova/gophkeeper/internal/server/api"
	"github.com/AlenaMolokova/gophkeeper/internal/server/auth"
	"github.com/AlenaMolokova/gophkeeper/internal/server/storage"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// startServer starts the gRPC server.
func startServer(grpcServer *grpc.Server, address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Server started on %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// logPoolStats logs pool statistics periodically.
// This function runs in a goroutine and logs connection pool statistics
// every 5 minutes to help monitor pool health and performance.
//
// The logged statistics include:
// - Total connections in the pool.
// - Idle connections available for use.
// - Acquired connections currently in use.
// - Connections being constructed.
func logPoolStats(storage *storage.PostgresStorage) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		stats := storage.GetPoolStats()
		log.Printf("Pool stats - Total: %d, Idle: %d, Acquired: %d, Constructing: %d",
			stats.TotalConns, stats.IdleConns, stats.AcquiredConns, stats.Constructing)
	}
}

// main starts the GophKeeper server.
// The function initializes the database connection with configurable pooling,
// sets up authentication, creates the gRPC server with TLS, and starts listening
// for client connections.
//
// Required environment variables:
//   - DATABASE_DSN: PostgreSQL connection string
//   - JWT_SECRET: Secret key for JWT token signing
//
// Optional environment variables for pool configuration:
//   - DB_MAX_CONNS, DB_MIN_CONNS, DB_MAX_CONN_LIFETIME, etc.
func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Load TLS credentials
	creds, err := credentials.NewServerTLSFromFile(cfg.TLS.CertFile, cfg.TLS.KeyFile)
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	// Create pool configuration from config
	poolConfig := &storage.PoolConfig{
		MaxConns:          cfg.Database.MaxConns,
		MinConns:          cfg.Database.MinConns,
		MaxConnLifetime:   cfg.Database.MaxConnLifetime,
		MaxConnIdleTime:   cfg.Database.MaxConnIdleTime,
		HealthCheckPeriod: cfg.Database.HealthCheckPeriod,
	}

	log.Printf("Database pool config - MaxConns: %d, MinConns: %d, MaxLifetime: %v, MaxIdleTime: %v, HealthCheckPeriod: %v",
		poolConfig.MaxConns, poolConfig.MinConns, poolConfig.MaxConnLifetime, poolConfig.MaxConnIdleTime, poolConfig.HealthCheckPeriod)

	// Initialize storage with custom pool configuration
	storage, err := storage.NewPostgresStorageWithConfig(ctx, cfg.GetDSN(), poolConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer storage.Close(ctx)

	// Start pool statistics logging
	go logPoolStats(storage)

	// Initialize auth
	authService := auth.NewAuth(storage, cfg.JWT.Secret, cfg.JWT.TTL)

	// Initialize gRPC servers
	userServer := api.NewUserServer(authService, cfg.GetJWTSecret())
	dataServer := api.NewDataServer(authService, storage, cfg.GetJWTSecret())
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// Register gRPC services
	clientapi.RegisterUserServiceServer(grpcServer, userServer)
	clientapi.RegisterDataServiceServer(grpcServer, dataServer)

	// Start server
	startServer(grpcServer, cfg.GetServerAddress())
}
