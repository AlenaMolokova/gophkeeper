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
	"os"
	"strconv"
	"time"

	"github.com/AlenaMolokova/gophkeeper/internal/server/api"
	"github.com/AlenaMolokova/gophkeeper/internal/server/auth"
	"github.com/AlenaMolokova/gophkeeper/internal/server/storage"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// getPoolConfig creates pool configuration from environment variables.
// It reads database connection pool settings from environment variables
// and returns a PoolConfig struct with the specified values.
//
// Supported environment variables:
//   - DB_MAX_CONNS: Maximum number of connections (default: 20)
//   - DB_MIN_CONNS: Minimum number of connections (default: 5)
//   - DB_MAX_CONN_LIFETIME: Connection lifetime (default: 1h)
//   - DB_MAX_CONN_IDLE_TIME: Idle connection timeout (default: 30m)
//   - DB_HEALTH_CHECK_PERIOD: Health check interval (default: 1m)
//
// If environment variables are not set or invalid, default values are used.
func getPoolConfig() *storage.PoolConfig {
	config := storage.DefaultPoolConfig()

	if maxConns := os.Getenv("DB_MAX_CONNS"); maxConns != "" {
		if val, err := strconv.ParseInt(maxConns, 10, 32); err == nil {
			config.MaxConns = int32(val)
		}
	}

	if minConns := os.Getenv("DB_MIN_CONNS"); minConns != "" {
		if val, err := strconv.ParseInt(minConns, 10, 32); err == nil {
			config.MinConns = int32(val)
		}
	}

	if maxLifetime := os.Getenv("DB_MAX_CONN_LIFETIME"); maxLifetime != "" {
		if val, err := time.ParseDuration(maxLifetime); err == nil {
			config.MaxConnLifetime = val
		}
	}

	if maxIdleTime := os.Getenv("DB_MAX_CONN_IDLE_TIME"); maxIdleTime != "" {
		if val, err := time.ParseDuration(maxIdleTime); err == nil {
			config.MaxConnIdleTime = val
		}
	}

	if healthCheckPeriod := os.Getenv("DB_HEALTH_CHECK_PERIOD"); healthCheckPeriod != "" {
		if val, err := time.ParseDuration(healthCheckPeriod); err == nil {
			config.HealthCheckPeriod = val
		}
	}

	return config
}

// startServer starts the gRPC server.
func startServer(grpcServer *grpc.Server) {
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Server started on :50051")
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
func logPoolStats(storage storage.Storage) {
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

	// Check required environment variables first
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN environment variable is not set")
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// Load TLS credentials before creating storage
	creds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	// Initialize storage with custom pool configuration
	poolConfig := getPoolConfig()
	log.Printf("Database pool config - MaxConns: %d, MinConns: %d, MaxLifetime: %v, MaxIdleTime: %v, HealthCheckPeriod: %v",
		poolConfig.MaxConns, poolConfig.MinConns, poolConfig.MaxConnLifetime, poolConfig.MaxConnIdleTime, poolConfig.HealthCheckPeriod)

	storage, err := storage.NewPostgresStorageWithConfig(ctx, dsn, poolConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer storage.Close(ctx)

	// Start pool statistics logging
	go logPoolStats(storage)

	// Initialize auth
	authService := auth.NewAuth(storage, string(jwtSecret), 24*3600*time.Second)

	// Initialize gRPC server
	server := api.NewServer(authService, storage, jwtSecret)
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// Register gRPC service
	clientapi.RegisterGophKeeperServer(grpcServer, server)

	// Start server
	startServer(grpcServer)
}
