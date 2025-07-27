package tui

import (
	"context"
	"testing"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// MockUserServiceClient is a mock implementation for testing.
type MockUserServiceClient struct {
	registerFunc func(ctx context.Context, req *clientapi.RegisterRequest, opts ...grpc.CallOption) (*clientapi.RegisterResponse, error)
	loginFunc    func(ctx context.Context, req *clientapi.LoginRequest, opts ...grpc.CallOption) (*clientapi.LoginResponse, error)
}

func (m *MockUserServiceClient) Register(ctx context.Context, req *clientapi.RegisterRequest, opts ...grpc.CallOption) (*clientapi.RegisterResponse, error) {
	return m.registerFunc(ctx, req, opts...)
}

func (m *MockUserServiceClient) Login(ctx context.Context, req *clientapi.LoginRequest, opts ...grpc.CallOption) (*clientapi.LoginResponse, error) {
	return m.loginFunc(ctx, req, opts...)
}

// MockDataServiceClient is a mock implementation for testing.
type MockDataServiceClient struct {
	addDataFunc    func(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error)
	getDataFunc    func(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error)
	editDataFunc   func(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error)
	deleteDataFunc func(ctx context.Context, req *clientapi.DeleteDataRequest, opts ...grpc.CallOption) (*clientapi.DeleteDataResponse, error)
}

func (m *MockDataServiceClient) AddData(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error) {
	return m.addDataFunc(ctx, req, opts...)
}

func (m *MockDataServiceClient) GetData(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error) {
	return m.getDataFunc(ctx, req, opts...)
}

func (m *MockDataServiceClient) EditData(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error) {
	return m.editDataFunc(ctx, req, opts...)
}

func (m *MockDataServiceClient) DeleteData(ctx context.Context, req *clientapi.DeleteDataRequest, opts ...grpc.CallOption) (*clientapi.DeleteDataResponse, error) {
	return m.deleteDataFunc(ctx, req, opts...)
}

func TestTUIAppCreation(t *testing.T) {
	mockUserClient := &MockUserServiceClient{}
	mockDataClient := &MockDataServiceClient{}

	app := &TUIApp{
		userClient: mockUserClient,
		dataClient: mockDataClient,
	}

	require.NotNil(t, app)
	assert.Equal(t, mockUserClient, app.userClient)
	assert.Equal(t, mockDataClient, app.dataClient)
}

func TestHandleRegister(t *testing.T) {
	mockUserClient := &MockUserServiceClient{
		registerFunc: func(ctx context.Context, req *clientapi.RegisterRequest, opts ...grpc.CallOption) (*clientapi.RegisterResponse, error) {
			return &clientapi.RegisterResponse{Token: "test-jwt-token"}, nil
		},
	}

	app := &TUIApp{
		userClient: mockUserClient,
		app:        tview.NewApplication(),
	}

	err := app.handleRegister("test@example.com", "password123")
	require.NoError(t, err)
}

func TestHandleRegisterError(t *testing.T) {
	mockUserClient := &MockUserServiceClient{
		registerFunc: func(ctx context.Context, req *clientapi.RegisterRequest, opts ...grpc.CallOption) (*clientapi.RegisterResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		userClient: mockUserClient,
		app:        tview.NewApplication(),
	}

	err := app.handleRegister("test@example.com", "password123")
	assert.Error(t, err)
}

func TestHandleLogin(t *testing.T) {
	mockUserClient := &MockUserServiceClient{
		loginFunc: func(ctx context.Context, req *clientapi.LoginRequest, opts ...grpc.CallOption) (*clientapi.LoginResponse, error) {
			return &clientapi.LoginResponse{Token: "test-jwt-token"}, nil
		},
	}

	app := &TUIApp{
		userClient: mockUserClient,
		app:        tview.NewApplication(),
	}

	err := app.handleLogin("test@example.com", "password123")
	require.NoError(t, err)
}

func TestHandleLoginError(t *testing.T) {
	mockUserClient := &MockUserServiceClient{
		loginFunc: func(ctx context.Context, req *clientapi.LoginRequest, opts ...grpc.CallOption) (*clientapi.LoginResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		userClient: mockUserClient,
		app:        tview.NewApplication(),
	}

	err := app.handleLogin("test@example.com", "password123")
	assert.Error(t, err)
}

func TestHandleAddData(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		addDataFunc: func(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, req.Data.Type)
			return &clientapi.AddDataResponse{Id: "data-id"}, nil
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	values := map[string]string{
		"JWT":     "jwt-token",
		"Type":    "login",
		"Payload": "test-payload",
	}

	err := app.handleAddData(values)
	require.NoError(t, err)
}

func TestHandleAddDataError(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		addDataFunc: func(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	values := map[string]string{
		"JWT":     "jwt-token",
		"Type":    "login",
		"Payload": "test-payload",
	}

	err := app.handleAddData(values)
	assert.Error(t, err)
}

func TestHandleGetData(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		getDataFunc: func(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error) {
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

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	values := map[string]string{
		"JWT": "jwt-token",
		"ID":  "data-id",
	}

	err := app.handleGetData(values)
	require.NoError(t, err)
}

func TestHandleGetDataError(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		getDataFunc: func(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	values := map[string]string{
		"JWT": "jwt-token",
		"ID":  "data-id",
	}

	err := app.handleGetData(values)
	assert.Error(t, err)
}

func TestHandleEditData(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		editDataFunc: func(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.Equal(t, "data-id", req.Data.Id)
			assert.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, req.Data.Type)
			return &clientapi.EditDataResponse{Id: "data-id"}, nil
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	values := map[string]string{
		"JWT":     "jwt-token",
		"ID":      "data-id",
		"Type":    "login",
		"Payload": "new-payload",
	}

	err := app.handleEditData(values)
	require.NoError(t, err)
}

func TestHandleEditDataError(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		editDataFunc: func(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	values := map[string]string{
		"JWT":     "jwt-token",
		"ID":      "data-id",
		"Type":    "login",
		"Payload": "new-payload",
	}

	err := app.handleEditData(values)
	assert.Error(t, err)
}

func TestHandleDeleteData(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		deleteDataFunc: func(ctx context.Context, req *clientapi.DeleteDataRequest, opts ...grpc.CallOption) (*clientapi.DeleteDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.Equal(t, "data-id", req.Id)
			return &clientapi.DeleteDataResponse{}, nil
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	values := map[string]string{
		"JWT": "jwt-token",
		"ID":  "data-id",
	}

	err := app.handleDeleteData(values)
	require.NoError(t, err)
}

func TestHandleDeleteDataError(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		deleteDataFunc: func(ctx context.Context, req *clientapi.DeleteDataRequest, opts ...grpc.CallOption) (*clientapi.DeleteDataResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	values := map[string]string{
		"JWT": "jwt-token",
		"ID":  "data-id",
	}

	err := app.handleDeleteData(values)
	assert.Error(t, err)
}
