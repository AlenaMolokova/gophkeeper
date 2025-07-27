// Package applogic provides high-level business logic functions for the GophKeeper client.
// It encapsulates common operations like user authentication, data encryption/decryption,
// and gRPC communication with the server.
//
// The package serves as an abstraction layer between the UI and the underlying
// gRPC clients, providing a clean API for client applications.
package applogic

import (
	"context"
	"time"

	clientcrypto "github.com/AlenaMolokova/gophkeeper/internal/client/crypto"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// RegisterUser registers a new user and returns a JWT token.
// It sends a registration request to the server and returns the authentication token
// upon successful account creation.
//
// The function handles the complete registration flow and provides
// the token needed for subsequent authenticated requests.
func RegisterUser(ctx context.Context, userClient clientapi.UserServiceClient, email, password string) (string, error) {
	resp, err := userClient.Register(ctx, &clientapi.RegisterRequest{Email: email, Password: password})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

// LoginUser authenticates a user and returns a JWT token.
// It sends a login request to the server and returns the authentication token
// upon successful authentication.
//
// The function handles the complete authentication flow and provides
// the token needed for subsequent authenticated requests.
func LoginUser(ctx context.Context, userClient clientapi.UserServiceClient, email, password string) (string, error) {
	resp, err := userClient.Login(ctx, &clientapi.LoginRequest{Email: email, Password: password})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

// AddData adds new data with encrypted payload.
// It encrypts the provided payload using client-side encryption,
// converts the data type string to the appropriate enum value,
// and sends the encrypted data to the server.
//
// The function ensures that sensitive data is encrypted before transmission
// and returns the unique data ID for future reference.
func AddData(ctx context.Context, dataClient clientapi.DataServiceClient, token, dataType, payload string) (string, error) {
	encPayload, err := clientcrypto.EncryptData([]byte(payload))
	if err != nil {
		return "", err
	}

	// Convert string dataType to DataType enum
	var dt clientapi.DataType
	switch dataType {
	case "login":
		dt = clientapi.DataType_DATA_TYPE_LOGIN
	case "text":
		dt = clientapi.DataType_DATA_TYPE_TEXT
	case "binary":
		dt = clientapi.DataType_DATA_TYPE_BINARY
	case "card":
		dt = clientapi.DataType_DATA_TYPE_CARD
	case "otp":
		dt = clientapi.DataType_DATA_TYPE_OTP
	default:
		dt = clientapi.DataType_DATA_TYPE_UNSPECIFIED
	}

	resp, err := dataClient.AddData(ctx, &clientapi.AddDataRequest{
		Token: token,
		Data: &clientapi.Data{
			Type:      dt,
			Payload:   encPayload,
			Timestamp: time.Now().Unix(),
		},
	})
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

// GetData retrieves and decrypts data by ID.
// It fetches encrypted data from the server and decrypts it using
// client-side decryption before returning both the data structure
// and the decrypted payload.
//
// The function ensures that sensitive data is properly decrypted
// for local use while maintaining the original data structure.
func GetData(ctx context.Context, dataClient clientapi.DataServiceClient, token, id string) (*clientapi.Data, string, error) {
	resp, err := dataClient.GetData(ctx, &clientapi.GetDataRequest{Token: token, Id: id})
	if err != nil {
		return nil, "", err
	}
	decPayload, err := clientcrypto.DecryptData(resp.Data.Payload)
	if err != nil {
		return nil, "", err
	}
	return resp.Data, string(decPayload), nil
}

// EditData updates existing data with encrypted payload.
// It encrypts the provided payload, converts the data type string to enum,
// and sends the updated encrypted data to the server.
//
// The function ensures that sensitive data is encrypted before transmission
// and maintains data integrity by preserving the original data ID.
func EditData(ctx context.Context, dataClient clientapi.DataServiceClient, token, id, dataType, payload string) (string, error) {
	encPayload, err := clientcrypto.EncryptData([]byte(payload))
	if err != nil {
		return "", err
	}

	// Convert string dataType to DataType enum
	var dt clientapi.DataType
	switch dataType {
	case "login":
		dt = clientapi.DataType_DATA_TYPE_LOGIN
	case "text":
		dt = clientapi.DataType_DATA_TYPE_TEXT
	case "binary":
		dt = clientapi.DataType_DATA_TYPE_BINARY
	case "card":
		dt = clientapi.DataType_DATA_TYPE_CARD
	case "otp":
		dt = clientapi.DataType_DATA_TYPE_OTP
	default:
		dt = clientapi.DataType_DATA_TYPE_UNSPECIFIED
	}

	resp, err := dataClient.EditData(ctx, &clientapi.EditDataRequest{
		Token: token,
		Data: &clientapi.Data{
			Id:        id,
			Type:      dt,
			Payload:   encPayload,
			Timestamp: time.Now().Unix(),
		},
	})
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

// DeleteData removes data by ID from the server.
// It sends a deletion request to the server for the specified data ID,
// ensuring that the data is permanently removed from the server storage.
func DeleteData(ctx context.Context, dataClient clientapi.DataServiceClient, token, id string) error {
	_, err := dataClient.DeleteData(ctx, &clientapi.DeleteDataRequest{Token: token, Id: id})
	return err
}
