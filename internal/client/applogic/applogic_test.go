package applogic

import (
	"context"
	"testing"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// MockUserServiceClient is a mock implementation of the UserServiceClient for testing.
type MockUserServiceClient struct {
	registerFunc func(ctx context.Context, req *clientapi.RegisterRequest) (*clientapi.RegisterResponse, error)
	loginFunc    func(ctx context.Context, req *clientapi.LoginRequest) (*clientapi.LoginResponse, error)
}

func (m *MockUserServiceClient) Register(ctx context.Context, req *clientapi.RegisterRequest, opts ...grpc.CallOption) (*clientapi.RegisterResponse, error) {
	return m.registerFunc(ctx, req)
}

func (m *MockUserServiceClient) Login(ctx context.Context, req *clientapi.LoginRequest, opts ...grpc.CallOption) (*clientapi.LoginResponse, error) {
	return m.loginFunc(ctx, req)
}

// MockDataServiceClient is a mock implementation of the DataServiceClient for testing.
type MockDataServiceClient struct {
	addDataFunc    func(ctx context.Context, req *clientapi.AddDataRequest) (*clientapi.AddDataResponse, error)
	getDataFunc    func(ctx context.Context, req *clientapi.GetDataRequest) (*clientapi.GetDataResponse, error)
	editDataFunc   func(ctx context.Context, req *clientapi.EditDataRequest) (*clientapi.EditDataResponse, error)
	deleteDataFunc func(ctx context.Context, req *clientapi.DeleteDataRequest) (*clientapi.DeleteDataResponse, error)
}

func (m *MockDataServiceClient) AddData(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error) {
	return m.addDataFunc(ctx, req)
}

func (m *MockDataServiceClient) GetData(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error) {
	return m.getDataFunc(ctx, req)
}

func (m *MockDataServiceClient) EditData(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error) {
	return m.editDataFunc(ctx, req)
}

func (m *MockDataServiceClient) DeleteData(ctx context.Context, req *clientapi.DeleteDataRequest, opts ...grpc.CallOption) (*clientapi.DeleteDataResponse, error) {
	return m.deleteDataFunc(ctx, req)
}

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockUserServiceClient{
		registerFunc: func(ctx context.Context, req *clientapi.RegisterRequest) (*clientapi.RegisterResponse, error) {
			assert.Equal(t, "test@example.com", req.Email)
			assert.Equal(t, "password123", req.Password)
			return &clientapi.RegisterResponse{Token: "jwt-token-123"}, nil
		},
	}

	token, err := RegisterUser(ctx, mockClient, "test@example.com", "password123")
	require.NoError(t, err)
	assert.Equal(t, "jwt-token-123", token)
}

func TestLoginUser(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockUserServiceClient{
		loginFunc: func(ctx context.Context, req *clientapi.LoginRequest) (*clientapi.LoginResponse, error) {
			assert.Equal(t, "test@example.com", req.Email)
			assert.Equal(t, "password123", req.Password)
			return &clientapi.LoginResponse{Token: "jwt-token-456"}, nil
		},
	}

	token, err := LoginUser(ctx, mockClient, "test@example.com", "password123")
	require.NoError(t, err)
	assert.Equal(t, "jwt-token-456", token)
}

func TestAddData(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockDataServiceClient{
		addDataFunc: func(ctx context.Context, req *clientapi.AddDataRequest) (*clientapi.AddDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, req.Data.Type)
			assert.NotEmpty(t, req.Data.Payload) // Should be encrypted
			return &clientapi.AddDataResponse{Id: "data-id-123"}, nil
		},
	}

	id, err := AddData(ctx, mockClient, "jwt-token", "login", "test-payload")
	require.NoError(t, err)
	assert.Equal(t, "data-id-123", id)
}

func TestGetData(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockDataServiceClient{
		getDataFunc: func(ctx context.Context, req *clientapi.GetDataRequest) (*clientapi.GetDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.Equal(t, "data-id", req.Id)
			return &clientapi.GetDataResponse{
				Data: &clientapi.Data{
					Id:      "data-id",
					Type:    clientapi.DataType_DATA_TYPE_LOGIN,
					Payload: []byte("encrypted-payload"),
				},
			}, nil
		},
	}

	data, decPayload, err := GetData(ctx, mockClient, "jwt-token", "data-id")
	require.NoError(t, err)
	assert.Equal(t, "data-id", data.Id)
	assert.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, data.Type)
	assert.NotEmpty(t, decPayload)
}

func TestEditData(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockDataServiceClient{
		editDataFunc: func(ctx context.Context, req *clientapi.EditDataRequest) (*clientapi.EditDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.Equal(t, "data-id", req.Data.Id)
			assert.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, req.Data.Type)
			assert.NotEmpty(t, req.Data.Payload) // Should be encrypted
			return &clientapi.EditDataResponse{Id: "data-id"}, nil
		},
	}

	id, err := EditData(ctx, mockClient, "jwt-token", "data-id", "login", "new-payload")
	require.NoError(t, err)
	assert.Equal(t, "data-id", id)
}

func TestDeleteData(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockDataServiceClient{
		deleteDataFunc: func(ctx context.Context, req *clientapi.DeleteDataRequest) (*clientapi.DeleteDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.Equal(t, "data-id", req.Id)
			return &clientapi.DeleteDataResponse{}, nil
		},
	}

	err := DeleteData(ctx, mockClient, "jwt-token", "data-id")
	require.NoError(t, err)
}
