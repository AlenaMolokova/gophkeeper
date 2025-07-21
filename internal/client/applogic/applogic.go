package applogic

import (
	"context"
	"time"

	clientcrypto "github.com/AlenaMolokova/gophkeeper/internal/client/crypto"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// RegisterUser registers a new user and returns a JWT token.
func RegisterUser(ctx context.Context, client clientapi.GophKeeperClient, email, password string) (string, error) {
	resp, err := client.Register(ctx, &clientapi.RegisterRequest{Email: email, Password: password})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

// LoginUser authenticates a user and returns a JWT token.
func LoginUser(ctx context.Context, client clientapi.GophKeeperClient, email, password string) (string, error) {
	resp, err := client.Login(ctx, &clientapi.LoginRequest{Email: email, Password: password})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

// AddData adds new data with encrypted payload.
func AddData(ctx context.Context, client clientapi.GophKeeperClient, token, dataType, payload string) (string, error) {
	encPayload, err := clientcrypto.EncryptData([]byte(payload))
	if err != nil {
		return "", err
	}
	resp, err := client.AddData(ctx, &clientapi.AddDataRequest{
		Token: token,
		Data: &clientapi.Data{
			Type:      dataType,
			Payload:   encPayload,
			Metadata:  map[string]string{},
			Timestamp: time.Now().Unix(),
		},
	})
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

// GetData retrieves and decrypts data by ID.
func GetData(ctx context.Context, client clientapi.GophKeeperClient, token, id string) (*clientapi.Data, string, error) {
	resp, err := client.GetData(ctx, &clientapi.GetDataRequest{Token: token, Id: id})
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
func EditData(ctx context.Context, client clientapi.GophKeeperClient, token, id, dataType, payload string) (string, error) {
	encPayload, err := clientcrypto.EncryptData([]byte(payload))
	if err != nil {
		return "", err
	}
	resp, err := client.EditData(ctx, &clientapi.EditDataRequest{
		Token: token,
		Data: &clientapi.Data{
			Id:        id,
			Type:      dataType,
			Payload:   encPayload,
			Metadata:  map[string]string{},
			Timestamp: time.Now().Unix(),
		},
	})
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

// DeleteData removes data by ID.
func DeleteData(ctx context.Context, client clientapi.GophKeeperClient, token, id string) error {
	_, err := client.DeleteData(ctx, &clientapi.DeleteDataRequest{Token: token, Id: id})
	return err
}
