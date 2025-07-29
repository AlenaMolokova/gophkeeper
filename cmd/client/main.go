// Package main provides the GophKeeper client application.
// The client implements a secure password manager with CLI and TUI interfaces,
// supporting encrypted data storage, synchronization, and secure communication.
//
// The client supports:
// - JWT-based authentication
// - Encrypted data storage and retrieval
// - CLI and TUI interfaces
// - TLS-secured gRPC communication
// - Local data synchronization
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	clientapplogic "github.com/AlenaMolokova/gophkeeper/internal/client/applogic"
	clientauth "github.com/AlenaMolokova/gophkeeper/internal/client/auth"
	clientcrypto "github.com/AlenaMolokova/gophkeeper/internal/client/crypto"
	clientsync "github.com/AlenaMolokova/gophkeeper/internal/client/sync"
	clienttui "github.com/AlenaMolokova/gophkeeper/internal/client/tui"
	"github.com/AlenaMolokova/gophkeeper/internal/config"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

var (
	// Version holds the version of the application.
	Version = "dev"
	// BuildDate holds the build date of the application.
	BuildDate = "unknown"
	// certPath holds the path to the TLS certificate.
	certPath string
)

// init initializes global flags.
func init() {
	clientConfig := config.NewClientConfig()
	flag.StringVar(&certPath, "cert", clientConfig.CertPath, "Path to the server TLS certificate")
	flag.Parse()
}

// loadTLSConfig loads TLS configuration from the specified certificate path.
// It returns the TLS credentials or an error if the certificate cannot be loaded.
func loadTLSConfig(certPath string) (credentials.TransportCredentials, error) {
	if certPath == "" {
		return nil, fmt.Errorf("certificate path is empty")
	}
	// Ensure path is compatible with the operating system
	certPath = filepath.Clean(certPath)
	log.Printf("Loading TLS certificate from: %s", certPath)
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate from %s: %w", certPath, err)
	}
	return creds, nil
}

// ClientConnections holds both user and data service connections.
type ClientConnections struct {
	UserClient clientapi.UserServiceClient
	DataClient clientapi.DataServiceClient
	Conn       *grpc.ClientConn
}

// connectToServer establishes a connection to the gRPC server.
// It returns the clients and connection, or an error if connection fails.
func connectToServer() (*ClientConnections, error) {
	creds, err := loadTLSConfig(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	clients := &ClientConnections{
		UserClient: clientapi.NewUserServiceClient(conn),
		DataClient: clientapi.NewDataServiceClient(conn),
		Conn:       conn,
	}
	return clients, nil
}

// executeWithConnection executes a function with server connections.
// It handles connection setup, cleanup, and error handling properly.
func executeWithConnection(operation func(context.Context, *ClientConnections) error) error {
	clients, err := connectToServer()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer clients.Conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return operation(ctx, clients)
}

// handleAuthCommand handles register and login commands with minimal logic.
// It parses command line arguments and delegates business logic to use cases.
func handleAuthCommand(command string) error {
	cmd := flag.NewFlagSet(command, flag.ExitOnError)
	email := cmd.String("email", "", "User email")
	password := cmd.String("password", "", "User password")
	if err := cmd.Parse(os.Args[2:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if *email == "" || *password == "" {
		fmt.Printf("Example: gophkeeper %s --email user@example.com --password secret\n", command)
		return fmt.Errorf("email and password are required")
	}

	return executeWithConnection(func(ctx context.Context, clients *ClientConnections) error {
		var token string
		var err error

		switch command {
		case "register":
			registerUsecase := clientauth.NewRegisterUsecase(clients.UserClient)
			token, err = registerUsecase.Execute(ctx, *email, *password)
			if err != nil {
				return fmt.Errorf("registration failed: %w", err)
			}
			fmt.Printf("User successfully registered. JWT: %s\n", token)
		case "login":
			loginUsecase := clientauth.NewLoginUsecase(clients.UserClient)
			token, err = loginUsecase.Execute(ctx, *email, *password)
			if err != nil {
				return fmt.Errorf("login failed: %w", err)
			}
			fmt.Printf("Login successful. JWT: %s\n", token)
		}
		return nil
	})
}

// handleAddCommand handles the add data command.
func handleAddCommand() error {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	token := addCmd.String("token", "", "JWT token")
	type_ := addCmd.String("type", "", "Data type (login, text, binary, card, otp)")
	payload := addCmd.String("payload", "", "Data (string, will be converted to []byte)")
	if err := addCmd.Parse(os.Args[2:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if *token == "" || *type_ == "" || *payload == "" {
		fmt.Println("Example: gophkeeper add --token <JWT> --type login --payload mysecret")
		return fmt.Errorf("token, type, and payload are required")
	}

	return executeWithConnection(func(ctx context.Context, clients *ClientConnections) error {
		id, err := clientapplogic.AddData(ctx, clients.DataClient, *token, *type_, *payload)
		if err != nil {
			return fmt.Errorf("failed to add data: %w", err)
		}
		fmt.Printf("Data successfully added. ID: %s\n", id)
		return nil
	})
}

// handleGetCommand handles the get data command.
func handleGetCommand() error {
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	token := getCmd.String("token", "", "JWT token")
	id := getCmd.String("id", "", "Data ID")
	if err := getCmd.Parse(os.Args[2:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if *token == "" || *id == "" {
		fmt.Println("Example: gophkeeper get --token <JWT> --id <data_id>")
		return fmt.Errorf("token and id are required")
	}

	return executeWithConnection(func(ctx context.Context, clients *ClientConnections) error {
		data, decPayload, err := clientapplogic.GetData(ctx, clients.DataClient, *token, *id)
		if err != nil {
			return fmt.Errorf("failed to get data: %w", err)
		}
		fmt.Printf("Data: ID=%s, Type=%s, Payload=%s, Timestamp=%d\n",
			data.Id, data.Type, decPayload, data.Timestamp)
		return nil
	})
}

// handleEditCommand handles the edit data command.
func handleEditCommand() error {
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	token := editCmd.String("token", "", "JWT token")
	id := editCmd.String("id", "", "Data ID")
	type_ := editCmd.String("type", "", "Data type")
	payload := editCmd.String("payload", "", "Data (string, will be converted to []byte)")
	if err := editCmd.Parse(os.Args[2:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if *token == "" || *id == "" || *type_ == "" || *payload == "" {
		fmt.Println("Example: gophkeeper edit --token <JWT> --id <data_id> --type login --payload newsecret")
		return fmt.Errorf("token, id, type, and payload are required")
	}

	return executeWithConnection(func(ctx context.Context, clients *ClientConnections) error {
		newID, err := clientapplogic.EditData(ctx, clients.DataClient, *token, *id, *type_, *payload)
		if err != nil {
			return fmt.Errorf("failed to edit data: %w", err)
		}
		fmt.Printf("Data successfully updated. ID: %s\n", newID)
		return nil
	})
}

// handleDeleteCommand handles the delete data command.
func handleDeleteCommand() error {
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	token := deleteCmd.String("token", "", "JWT token")
	id := deleteCmd.String("id", "", "Data ID")
	if err := deleteCmd.Parse(os.Args[2:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if *token == "" || *id == "" {
		fmt.Println("Example: gophkeeper delete --token <JWT> --id <data_id>")
		return fmt.Errorf("token and id are required")
	}

	return executeWithConnection(func(ctx context.Context, clients *ClientConnections) error {
		if err := clientapplogic.DeleteData(ctx, clients.DataClient, *token, *id); err != nil {
			return fmt.Errorf("failed to delete data: %w", err)
		}
		fmt.Println("Data successfully deleted.")
		return nil
	})
}

// handleInitKeyCommand handles the key initialization command.
func handleInitKeyCommand() error {
	key, err := clientcrypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}
	err = clientcrypto.SaveKey(key)
	if err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}
	fmt.Println("Key successfully generated and saved.")
	return nil
}

// handleTUICommand handles the TUI command.
func handleTUICommand() error {
	return clienttui.RunTUI()
}

// handleSyncUploadCommand handles the sync upload command.
func handleSyncUploadCommand() error {
	syncer, err := clientsync.NewDefaultSyncer()
	if err != nil {
		return fmt.Errorf("failed to initialize syncer: %w", err)
	}
	if err := syncer.Upload(context.Background()); err != nil {
		return fmt.Errorf("failed to upload data: %w", err)
	}
	fmt.Println("Data upload to server completed.")
	return nil
}

// handleSyncDownloadCommand handles the sync download command.
func handleSyncDownloadCommand() error {
	syncer, err := clientsync.NewDefaultSyncer()
	if err != nil {
		return fmt.Errorf("failed to initialize syncer: %w", err)
	}
	if err := syncer.Download(context.Background()); err != nil {
		return fmt.Errorf("failed to download data: %w", err)
	}
	fmt.Println("Data download from server completed.")
	return nil
}

// showUsage displays the usage information.
func showUsage() {
	fmt.Println("Usage: gophkeeper <command> [parameters]")
	fmt.Println("Available commands: register, login, add, get, edit, delete, version, tui, sync-upload, sync-download, init-key")
	fmt.Println("Use --cert to specify TLS certificate path (or set GOPHKEEPER_CERT environment variable)")
}

// showVersion displays the version information.
func showVersion() {
	fmt.Printf("GophKeeper CLI\nVersion: %s\nBuild Date: %s\n", Version, BuildDate)
}

// main is the entry point for the GophKeeper CLI client.
func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "register", "login":
		err = handleAuthCommand(os.Args[1])
	case "init-key":
		err = handleInitKeyCommand()
	case "add":
		err = handleAddCommand()
	case "get":
		err = handleGetCommand()
	case "edit":
		err = handleEditCommand()
	case "delete":
		err = handleDeleteCommand()
	case "version":
		showVersion()
		return
	case "tui":
		err = handleTUICommand()
	case "sync-upload":
		err = handleSyncUploadCommand()
	case "sync-download":
		err = handleSyncDownloadCommand()
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}
