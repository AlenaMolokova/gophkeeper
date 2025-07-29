package applogic

import (
	"context"
	"testing"

	"github.com/AlenaMolokova/gophkeeper/internal/client/applogic/mocks"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockUserServiceClientInterface(ctrl)

	mockClient.EXPECT().
		Register(gomock.Any(), &clientapi.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}).
		Return(&clientapi.RegisterResponse{Token: "jwt-token-123"}, nil)

	token, err := RegisterUser(context.Background(), mockClient, "test@example.com", "password123")
	require.NoError(t, err)
	assert.Equal(t, "jwt-token-123", token)
}

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockUserServiceClientInterface(ctrl)

	mockClient.EXPECT().
		Login(gomock.Any(), &clientapi.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}).
		Return(&clientapi.LoginResponse{Token: "jwt-token-456"}, nil)

	token, err := LoginUser(context.Background(), mockClient, "test@example.com", "password123")
	require.NoError(t, err)
	assert.Equal(t, "jwt-token-456", token)
}

func TestAddData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockDataServiceClientInterface(ctrl)

	mockClient.EXPECT().
		AddData(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, req *clientapi.AddDataRequest, opts ...interface{}) (*clientapi.AddDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, req.Data.Type)
			assert.NotEmpty(t, req.Data.Payload) // Should be encrypted
			return &clientapi.AddDataResponse{Id: "data-id-123"}, nil
		})

	id, err := AddData(context.Background(), mockClient, "jwt-token", "login", "test-payload")
	require.NoError(t, err)
	assert.Equal(t, "data-id-123", id)
}

func TestGetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockDataServiceClientInterface(ctrl)

	mockClient.EXPECT().
		GetData(gomock.Any(), &clientapi.GetDataRequest{
			Token: "jwt-token",
			Id:    "data-id",
		}).
		Return(&clientapi.GetDataResponse{
			Data: &clientapi.Data{
				Id:      "data-id",
				Type:    clientapi.DataType_DATA_TYPE_LOGIN,
				Payload: []byte("encrypted-payload"),
			},
		}, nil)

	data, decPayload, err := GetData(context.Background(), mockClient, "jwt-token", "data-id")
	require.NoError(t, err)
	assert.Equal(t, "data-id", data.Id)
	assert.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, data.Type)
	assert.NotEmpty(t, decPayload)
}

func TestEditData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockDataServiceClientInterface(ctrl)

	mockClient.EXPECT().
		EditData(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, req *clientapi.EditDataRequest, opts ...interface{}) (*clientapi.EditDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.Equal(t, "data-id", req.Data.Id)
			assert.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, req.Data.Type)
			assert.NotEmpty(t, req.Data.Payload) // Should be encrypted
			return &clientapi.EditDataResponse{Id: "data-id"}, nil
		})

	id, err := EditData(context.Background(), mockClient, "jwt-token", "data-id", "login", "new-payload")
	require.NoError(t, err)
	assert.Equal(t, "data-id", id)
}

func TestDeleteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockDataServiceClientInterface(ctrl)

	mockClient.EXPECT().
		DeleteData(gomock.Any(), &clientapi.DeleteDataRequest{
			Token: "jwt-token",
			Id:    "data-id",
		}).
		Return(&clientapi.DeleteDataResponse{}, nil)

	err := DeleteData(context.Background(), mockClient, "jwt-token", "data-id")
	require.NoError(t, err)
}
