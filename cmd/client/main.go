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
	clientcrypto "github.com/AlenaMolokova/gophkeeper/internal/client/crypto"
	clientsync "github.com/AlenaMolokova/gophkeeper/internal/client/sync"
	clienttui "github.com/AlenaMolokova/gophkeeper/internal/client/tui"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

var (
	// Version holds the current version of the application.
	Version = "dev"
	// BuildDate holds the build date of the application.
	BuildDate = "unknown"
	// certPath holds the path to the TLS certificate.
	certPath string
)

// init initializes global flags.
func init() {
	flag.StringVar(&certPath, "cert", os.Getenv("GOPHKEEPER_CERT"), "Path to the server TLS certificate")
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

// connectToServer establishes a connection to the gRPC server.
// It returns the client and connection, or an error if connection fails.
func connectToServer() (clientapi.GophKeeperClient, *grpc.ClientConn, error) {
	creds, err := loadTLSConfig(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := clientapi.NewGophKeeperClient(conn)
	return client, conn, nil
}

// executeWithConnection executes a function with a server connection.
// It handles connection setup, cleanup, and error handling properly.
func executeWithConnection(operation func(context.Context, clientapi.GophKeeperClient) error) {
	client, conn, err := connectToServer()
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := operation(ctx, client); err != nil {
		log.Printf("Operation failed: %v", err)
		os.Exit(1)
	}
}



// handleAuthCommand handles register and login commands with common logic.
func handleAuthCommand(command string) {
	cmd := flag.NewFlagSet(command, flag.ExitOnError)
	email := cmd.String("email", "", "User email")
	password := cmd.String("password", "", "User password")
	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	if *email == "" || *password == "" {
		fmt.Printf("Example: gophkeeper %s --email user@example.com --password secret\n", command)
		os.Exit(1)
	}

	executeWithConnection(func(ctx context.Context, client clientapi.GophKeeperClient) error {
		var token string
		var err error
		switch command {
		case "register":
			token, err = clientapplogic.RegisterUser(ctx, client, *email, *password)
			if err != nil {
				return fmt.Errorf("registration failed: %w", err)
			}
			fmt.Printf("User successfully registered. JWT: %s\n", token)
		case "login":
			token, err = clientapplogic.LoginUser(ctx, client, *email, *password)
			if err != nil {
				return fmt.Errorf("login failed: %w", err)
			}
			fmt.Printf("Login successful. JWT: %s\n", token)
		}
		return nil
	})
}

// main is the entry point for the GophKeeper CLI client.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gophkeeper <command> [parameters]")
		fmt.Println("Available commands: register, login, add, get, edit, delete, version, tui, sync-upload, sync-download, init-key")
		fmt.Println("Use --cert to specify TLS certificate path (or set GOPHKEEPER_CERT environment variable)")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "register", "login":
		handleAuthCommand(os.Args[1])

	case "init-key":
		key, err := clientcrypto.GenerateKey()
		if err != nil {
			log.Fatalf("Failed to generate key: %v", err)
		}
		err = clientcrypto.SaveKey(key)
		if err != nil {
			log.Fatalf("Failed to save key: %v", err)
		}
		fmt.Println("Key successfully generated and saved.")
		os.Exit(0)

	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		token := addCmd.String("token", "", "JWT token")
		type_ := addCmd.String("type", "", "Data type (login, text, binary, card, otp)")
		payload := addCmd.String("payload", "", "Data (string, will be converted to []byte)")
		if err := addCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("Failed to parse flags: %v", err)
		}

		if *token == "" || *type_ == "" || *payload == "" {
			fmt.Println("Example: gophkeeper add --token <JWT> --type login --payload mysecret")
			os.Exit(1)
		}

		executeWithConnection(func(ctx context.Context, client clientapi.GophKeeperClient) error {
			id, err := clientapplogic.AddData(ctx, client, *token, *type_, *payload)
			if err != nil {
				return fmt.Errorf("failed to add data: %w", err)
			}
			fmt.Printf("Data successfully added. ID: %s\n", id)
			return nil
		})

	case "get":
		getCmd := flag.NewFlagSet("get", flag.ExitOnError)
		token := getCmd.String("token", "", "JWT token")
		id := getCmd.String("id", "", "Data ID")
		if err := getCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("Failed to parse flags: %v", err)
		}

		if *token == "" || *id == "" {
			fmt.Println("Example: gophkeeper get --token <JWT> --id <data_id>")
			os.Exit(1)
		}

		executeWithConnection(func(ctx context.Context, client clientapi.GophKeeperClient) error {
			data, decPayload, err := clientapplogic.GetData(ctx, client, *token, *id)
			if err != nil {
				return fmt.Errorf("failed to get data: %w", err)
			}
			fmt.Printf("Data: ID=%s, Type=%s, Payload=%s, Metadata=%v, Timestamp=%d\n",
				data.Id, data.Type, decPayload, data.Metadata, data.Timestamp)
			return nil
		})

	case "edit":
		editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
		token := editCmd.String("token", "", "JWT token")
		id := editCmd.String("id", "", "Data ID")
		type_ := editCmd.String("type", "", "Data type")
		payload := editCmd.String("payload", "", "Data (string, will be converted to []byte)")
		if err := editCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("Failed to parse flags: %v", err)
		}

		if *token == "" || *id == "" || *type_ == "" || *payload == "" {
			fmt.Println("Example: gophkeeper edit --token <JWT> --id <data_id> --type login --payload newsecret")
			os.Exit(1)
		}

		executeWithConnection(func(ctx context.Context, client clientapi.GophKeeperClient) error {
			newID, err := clientapplogic.EditData(ctx, client, *token, *id, *type_, *payload)
			if err != nil {
				return fmt.Errorf("failed to edit data: %w", err)
			}
			fmt.Printf("Data successfully updated. ID: %s\n", newID)
			return nil
		})

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		token := deleteCmd.String("token", "", "JWT token")
		id := deleteCmd.String("id", "", "Data ID")
		if err := deleteCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("Failed to parse flags: %v", err)
		}

		if *token == "" || *id == "" {
			fmt.Println("Example: gophkeeper delete --token <JWT> --id <data_id>")
			os.Exit(1)
		}

		executeWithConnection(func(ctx context.Context, client clientapi.GophKeeperClient) error {
			if err := clientapplogic.DeleteData(ctx, client, *token, *id); err != nil {
				return fmt.Errorf("failed to delete data: %w", err)
			}
			fmt.Println("Data successfully deleted.")
			return nil
		})

	case "version":
		fmt.Printf("GophKeeper CLI\nVersion: %s\nBuild Date: %s\n", Version, BuildDate)
		os.Exit(0)

	case "tui":
		if err := clienttui.RunTUI(); err != nil {
			log.Fatalf("TUI failed: %v", err)
		}
		os.Exit(0)

	case "sync-upload":
		syncer, err := clientsync.NewDefaultSyncer()
		if err != nil {
			log.Fatalf("Failed to initialize syncer: %v", err)
		}
		if err := syncer.Upload(context.Background()); err != nil {
			log.Fatalf("Failed to upload data: %v", err)
		}
		fmt.Println("Data upload to server completed.")
		os.Exit(0)

	case "sync-download":
		syncer, err := clientsync.NewDefaultSyncer()
		if err != nil {
			log.Fatalf("Failed to initialize syncer: %v", err)
		}
		if err := syncer.Download(context.Background()); err != nil {
			log.Fatalf("Failed to download data: %v", err)
		}
		fmt.Println("Data download from server completed.")
		os.Exit(0)

	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}
