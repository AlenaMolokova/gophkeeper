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

// MockGophKeeperClient is a mock implementation for testing.
type MockGophKeeperClient struct {
	registerFunc   func(ctx context.Context, req *clientapi.RegisterRequest, opts ...grpc.CallOption) (*clientapi.RegisterResponse, error)
	loginFunc      func(ctx context.Context, req *clientapi.LoginRequest, opts ...grpc.CallOption) (*clientapi.LoginResponse, error)
	addDataFunc    func(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error)
	getDataFunc    func(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error)
	editDataFunc   func(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error)
	deleteDataFunc func(ctx context.Context, req *clientapi.DeleteDataRequest, opts ...grpc.CallOption) (*clientapi.DeleteDataResponse, error)
}

func (m *MockGophKeeperClient) Register(ctx context.Context, req *clientapi.RegisterRequest, opts ...grpc.CallOption) (*clientapi.RegisterResponse, error) {
	return m.registerFunc(ctx, req, opts...)
}

func (m *MockGophKeeperClient) Login(ctx context.Context, req *clientapi.LoginRequest, opts ...grpc.CallOption) (*clientapi.LoginResponse, error) {
	return m.loginFunc(ctx, req, opts...)
}

func (m *MockGophKeeperClient) AddData(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error) {
	return m.addDataFunc(ctx, req, opts...)
}

func (m *MockGophKeeperClient) GetData(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error) {
	return m.getDataFunc(ctx, req, opts...)
}

func (m *MockGophKeeperClient) EditData(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error) {
	return m.editDataFunc(ctx, req, opts...)
}

func (m *MockGophKeeperClient) DeleteData(ctx context.Context, req *clientapi.DeleteDataRequest, opts ...grpc.CallOption) (*clientapi.DeleteDataResponse, error) {
	return m.deleteDataFunc(ctx, req, opts...)
}

func TestTUIAppCreation(t *testing.T) {
	mockClient := &MockGophKeeperClient{}

	app := &TUIApp{
		client: mockClient,
	}

	require.NotNil(t, app)
	assert.Equal(t, mockClient, app.client)
}

func TestHandleRegister(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		registerFunc: func(ctx context.Context, req *clientapi.RegisterRequest, opts ...grpc.CallOption) (*clientapi.RegisterResponse, error) {
			return &clientapi.RegisterResponse{Token: "test-jwt-token"}, nil
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	err := app.handleRegister("test@example.com", "password123")
	require.NoError(t, err)
}

func TestHandleRegisterError(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		registerFunc: func(ctx context.Context, req *clientapi.RegisterRequest, opts ...grpc.CallOption) (*clientapi.RegisterResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	err := app.handleRegister("test@example.com", "password123")
	assert.Error(t, err)
}

func TestHandleLogin(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		loginFunc: func(ctx context.Context, req *clientapi.LoginRequest, opts ...grpc.CallOption) (*clientapi.LoginResponse, error) {
			return &clientapi.LoginResponse{Token: "test-jwt-token"}, nil
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	err := app.handleLogin("test@example.com", "password123")
	require.NoError(t, err)
}

func TestHandleLoginError(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		loginFunc: func(ctx context.Context, req *clientapi.LoginRequest, opts ...grpc.CallOption) (*clientapi.LoginResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	err := app.handleLogin("test@example.com", "password123")
	assert.Error(t, err)
}

func TestHandleAddData(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		addDataFunc: func(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error) {
			return &clientapi.AddDataResponse{Id: "test-data-id"}, nil
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	values := map[string]string{
		"JWT":     "test-jwt",
		"Type":    "login",
		"Payload": "test-payload",
	}

	err := app.handleAddData(values)
	require.NoError(t, err)
}

func TestHandleAddDataError(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		addDataFunc: func(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	values := map[string]string{
		"JWT":     "test-jwt",
		"Type":    "login",
		"Payload": "test-payload",
	}

	err := app.handleAddData(values)
	assert.Error(t, err)
}

func TestHandleGetData(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		getDataFunc: func(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error) {
			return &clientapi.GetDataResponse{
				Data: &clientapi.Data{
					Id:      "test-data-id",
					Type:    "login",
					Payload: []byte("encrypted-payload"),
				},
			}, nil
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	values := map[string]string{
		"JWT": "test-jwt",
		"ID":  "test-data-id",
	}

	err := app.handleGetData(values)
	require.NoError(t, err)
}

func TestHandleGetDataError(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		getDataFunc: func(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	values := map[string]string{
		"JWT": "test-jwt",
		"ID":  "test-data-id",
	}

	err := app.handleGetData(values)
	assert.Error(t, err)
}

func TestHandleEditData(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		editDataFunc: func(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error) {
			return &clientapi.EditDataResponse{Id: "test-data-id"}, nil
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	values := map[string]string{
		"JWT":     "test-jwt",
		"ID":      "test-data-id",
		"Type":    "login",
		"Payload": "updated-payload",
	}

	err := app.handleEditData(values)
	require.NoError(t, err)
}

func TestHandleEditDataError(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		editDataFunc: func(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	values := map[string]string{
		"JWT":     "test-jwt",
		"ID":      "test-data-id",
		"Type":    "login",
		"Payload": "updated-payload",
	}

	err := app.handleEditData(values)
	assert.Error(t, err)
}

func TestHandleDeleteData(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		deleteDataFunc: func(ctx context.Context, req *clientapi.DeleteDataRequest, opts ...grpc.CallOption) (*clientapi.DeleteDataResponse, error) {
			return &clientapi.DeleteDataResponse{}, nil
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	values := map[string]string{
		"JWT": "test-jwt",
		"ID":  "test-data-id",
	}

	err := app.handleDeleteData(values)
	require.NoError(t, err)
}

func TestHandleDeleteDataError(t *testing.T) {
	mockClient := &MockGophKeeperClient{
		deleteDataFunc: func(ctx context.Context, req *clientapi.DeleteDataRequest, opts ...grpc.CallOption) (*clientapi.DeleteDataResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		client: mockClient,
		app:    tview.NewApplication(),
	}

	values := map[string]string{
		"JWT": "test-jwt",
		"ID":  "test-data-id",
	}

	err := app.handleDeleteData(values)
	assert.Error(t, err)
}
