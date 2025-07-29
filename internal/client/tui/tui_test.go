package tui

import (
	"context"
	"fmt"
	"os"
	"testing"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/grpc"

	"github.com/nicksnyder/go-i18n/v2/i18n"
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

// Test helper functions for creating typed form data
func createTestAuthFormData() *AuthFormData {
	return &AuthFormData{
		Email:    "test@example.com",
		Password: "password123",
	}
}

func createTestJWTIDFormData() *JWTIDFormData {
	return &JWTIDFormData{
		JWT: "jwt-token",
		ID:  "data-id",
	}
}

func createTestAddDataFormData() *AddDataFormData {
	return &AddDataFormData{
		JWT:     "jwt-token",
		Type:    "login",
		Payload: "test-payload",
	}
}

func createTestEditDataFormData() *EditDataFormData {
	return &EditDataFormData{
		JWT:     "jwt-token",
		ID:      "data-id",
		Type:    "login",
		Payload: "updated-payload",
	}
}

// Helper functions for creating typed handlers from map[string]string handlers
func createJWTIDHandler(handler func(*JWTIDFormData) error) func(map[string]string) error {
	return func(values map[string]string) error {
		data := &JWTIDFormData{
			JWT: values["JWT"],
			ID:  values["ID"],
		}
		return handler(data)
	}
}

func createAddDataHandler(handler func(*AddDataFormData) error) func(map[string]string) error {
	return func(values map[string]string) error {
		data := &AddDataFormData{
			JWT:     values["JWT"],
			Type:    values["Type"],
			Payload: values["Payload"],
		}
		return handler(data)
	}
}

func createEditDataHandler(handler func(*EditDataFormData) error) func(map[string]string) error {
	return func(values map[string]string) error {
		data := &EditDataFormData{
			JWT:     values["JWT"],
			ID:      values["ID"],
			Type:    values["Type"],
			Payload: values["Payload"],
		}
		return handler(data)
	}
}

// TestTUIAppCreation tests TUIApp creation and initialization.
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

// TestShowModal tests the showModal method.
func TestShowModal(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Test that showModal doesn't panic
	assert.NotPanics(t, func() {
		app.showModal("Test message")
	})
}

// TestCreateAuthForm tests the createAuthForm method with different actions.
func TestCreateAuthForm(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	handler := func(data *AuthFormData) error {
		return nil
	}

	// Test register form
	registerForm := app.createAuthForm("Зарегистрироваться", handler)
	assert.NotNil(t, registerForm)

	// Test login form
	loginForm := app.createAuthForm("Войти", handler)
	assert.NotNil(t, loginForm)

	// Test unknown action (fallback)
	unknownForm := app.createAuthForm("Unknown", handler)
	assert.NotNil(t, unknownForm)
}

// TestCreateAuthFormWithError tests createAuthForm with error handling.
func TestCreateAuthFormWithError(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	errorHandler := func(data *AuthFormData) error {
		return assert.AnError
	}

	form := app.createAuthForm("Зарегистрироваться", errorHandler)
	assert.NotNil(t, form)
}

// TestCreateDataForm tests the createDataForm method with different actions.
func TestCreateDataForm(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Create typed handlers for each form type
	getHandler := func(data *JWTIDFormData) error {
		return nil
	}
	addHandler := func(data *AddDataFormData) error {
		return nil
	}
	editHandler := func(data *EditDataFormData) error {
		return nil
	}
	deleteHandler := func(data *JWTIDFormData) error {
		return nil
	}

	// Test get data form
	getForm := app.createDataForm("Получить", []string{"JWT", "ID"}, createJWTIDHandler(getHandler))
	assert.NotNil(t, getForm)

	// Test add data form
	addForm := app.createDataForm("Добавить", []string{"JWT", "Type", "Payload"}, createAddDataHandler(addHandler))
	assert.NotNil(t, addForm)

	// Test edit data form
	editForm := app.createDataForm("Изменить", []string{"JWT", "ID", "Type", "Payload"}, createEditDataHandler(editHandler))
	assert.NotNil(t, editForm)

	// Test delete data form
	deleteForm := app.createDataForm("Удалить", []string{"JWT", "ID"}, createJWTIDHandler(deleteHandler))
	assert.NotNil(t, deleteForm)

	// Test unknown action (fallback)
	unknownForm := app.createDataForm("Unknown", []string{"JWT", "ID"}, func(values map[string]string) error {
		return nil
	})
	assert.NotNil(t, unknownForm)
}

// TestCreateDataFormWithError tests createDataForm with error handling.
func TestCreateDataFormWithError(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	errorHandler := func(values map[string]string) error {
		return assert.AnError
	}

	form := app.createDataForm("Получить", []string{"JWT", "ID"}, errorHandler)
	assert.NotNil(t, form)
}

// TestHandleRegister tests user registration handling.
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

	err := app.handleRegister(createTestAuthFormData())
	require.NoError(t, err)
}

// TestHandleRegisterError tests user registration error handling.
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

	err := app.handleRegister(createTestAuthFormData())
	assert.Error(t, err)
}

// TestHandleLogin tests user login handling.
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

	err := app.handleLogin(createTestAuthFormData())
	require.NoError(t, err)
}

// TestHandleLoginError tests user login error handling.
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

	err := app.handleLogin(createTestAuthFormData())
	assert.Error(t, err)
}

// TestHandleAddData tests adding data handling.
func TestHandleAddData(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		addDataFunc: func(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.NotNil(t, req.Data)
			assert.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, req.Data.Type)
			// Payload is encrypted, so we can't check the exact string value
			assert.NotEmpty(t, req.Data.Payload)
			return &clientapi.AddDataResponse{Id: "new-data-id"}, nil
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	err := app.handleAddData(createTestAddDataFormData())
	require.NoError(t, err)
}

// TestHandleAddDataError tests adding data error handling.
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

	err := app.handleAddData(createTestAddDataFormData())
	assert.Error(t, err)
}

// TestHandleGetData tests retrieving data handling.
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

	err := app.handleGetData(createTestJWTIDFormData())
	require.NoError(t, err)
}

// TestHandleGetDataError tests retrieving data error handling.
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

	err := app.handleGetData(createTestJWTIDFormData())
	assert.Error(t, err)
}

// TestHandleEditData tests editing data handling.
func TestHandleEditData(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		editDataFunc: func(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error) {
			assert.Equal(t, "jwt-token", req.Token)
			assert.NotNil(t, req.Data)
			assert.Equal(t, "data-id", req.Data.Id)
			assert.Equal(t, clientapi.DataType_DATA_TYPE_LOGIN, req.Data.Type)
			// Payload is encrypted, so we can't check the exact string value
			assert.NotEmpty(t, req.Data.Payload)
			return &clientapi.EditDataResponse{Id: "data-id"}, nil
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	err := app.handleEditData(createTestEditDataFormData())
	require.NoError(t, err)
}

// TestHandleEditDataError tests editing data error handling.
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

	err := app.handleEditData(createTestEditDataFormData())
	assert.Error(t, err)
}

// TestHandleDeleteData tests deleting data handling.
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

	err := app.handleDeleteData(createTestJWTIDFormData())
	require.NoError(t, err)
}

// TestHandleDeleteDataError tests deleting data error handling.
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

	err := app.handleDeleteData(createTestJWTIDFormData())
	assert.Error(t, err)
}

// TestSetupMenu tests the setupMenu method.
func TestSetupMenu(t *testing.T) {
	// Create a proper localizer for testing
	localizer := &Localizer{
		bundle:    i18n.NewBundle(language.Russian),
		localizer: i18n.NewLocalizer(i18n.NewBundle(language.Russian), "ru"),
	}

	app := &TUIApp{
		app:       tview.NewApplication(),
		localizer: localizer,
	}

	// Test that setupMenu doesn't panic
	assert.NotPanics(t, func() {
		app.setupMenu()
	})

	// Test that menu is created
	assert.NotNil(t, app.menu)
}

// TestHandleAddDataWithInvalidType tests adding data with invalid type.
func TestHandleAddDataWithInvalidType(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		addDataFunc: func(ctx context.Context, req *clientapi.AddDataRequest, opts ...grpc.CallOption) (*clientapi.AddDataResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	// Test with invalid data type
	invalidData := &AddDataFormData{
		JWT:     "jwt-token",
		Type:    "invalid-type",
		Payload: "test-payload",
	}

	err := app.handleAddData(invalidData)
	assert.Error(t, err)
}

// TestHandleEditDataWithInvalidType tests editing data with invalid type.
func TestHandleEditDataWithInvalidType(t *testing.T) {
	mockDataClient := &MockDataServiceClient{
		editDataFunc: func(ctx context.Context, req *clientapi.EditDataRequest, opts ...grpc.CallOption) (*clientapi.EditDataResponse, error) {
			return nil, assert.AnError
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	// Test with invalid data type
	invalidData := &EditDataFormData{
		JWT:     "jwt-token",
		ID:      "data-id",
		Type:    "invalid-type",
		Payload: "test-payload",
	}

	err := app.handleEditData(invalidData)
	assert.Error(t, err)
}

// TestHandleAddDataWithEmptyValues tests adding data with empty values.
func TestHandleAddDataWithEmptyValues(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Test with empty values
	emptyData := &AddDataFormData{
		JWT:     "",
		Type:    "",
		Payload: "",
	}

	err := app.handleAddData(emptyData)
	assert.Error(t, err)
}

// TestHandleGetDataWithEmptyValues tests getting data with empty values.
func TestHandleGetDataWithEmptyValues(t *testing.T) {
	// Create a mock data client that returns error for empty values
	mockDataClient := &MockDataServiceClient{
		getDataFunc: func(ctx context.Context, req *clientapi.GetDataRequest, opts ...grpc.CallOption) (*clientapi.GetDataResponse, error) {
			return nil, fmt.Errorf("invalid request")
		},
	}

	app := &TUIApp{
		dataClient: mockDataClient,
		app:        tview.NewApplication(),
	}

	// Test with empty values
	emptyData := &JWTIDFormData{
		JWT: "",
		ID:  "",
	}

	err := app.handleGetData(emptyData)
	assert.Error(t, err)
}

// TestHandleEditDataWithEmptyValues tests editing data with empty values.
func TestHandleEditDataWithEmptyValues(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Test with empty values
	emptyData := &EditDataFormData{
		JWT:     "",
		ID:      "",
		Type:    "",
		Payload: "",
	}

	err := app.handleEditData(emptyData)
	assert.Error(t, err)
}

// TestHandleDeleteDataWithEmptyValues tests deleting data with empty values.
func TestHandleDeleteDataWithEmptyValues(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Test with empty values
	emptyData := &JWTIDFormData{
		JWT: "",
		ID:  "",
	}

	err := app.handleDeleteData(emptyData)
	assert.Error(t, err)
}

// TestHandleRegisterWithEmptyCredentials tests registration with empty credentials.
func TestHandleRegisterWithEmptyCredentials(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Test with empty credentials
	emptyData := &AuthFormData{
		Email:    "",
		Password: "",
	}

	err := app.handleRegister(emptyData)
	assert.Error(t, err)
}

// TestHandleLoginWithEmptyCredentials tests login with empty credentials.
func TestHandleLoginWithEmptyCredentials(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Test with empty credentials
	emptyData := &AuthFormData{
		Email:    "",
		Password: "",
	}

	err := app.handleLogin(emptyData)
	assert.Error(t, err)
}

// TestCreateAuthFormWithUnknownAction tests createAuthForm with unknown action.
func TestCreateAuthFormWithUnknownAction(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	handler := func(data *AuthFormData) error {
		return nil
	}

	form := app.createAuthForm("UnknownAction", handler)
	assert.NotNil(t, form)

	// Check that form has expected fields
	emailField := form.GetFormItemByLabel("Email")
	assert.NotNil(t, emailField)

	passwordField := form.GetFormItemByLabel("Пароль")
	assert.NotNil(t, passwordField)
}

// TestCreateDataFormWithUnknownAction tests createDataForm with unknown action.
func TestCreateDataFormWithUnknownAction(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	handler := func(values map[string]string) error {
		return nil
	}

	form := app.createDataForm("UnknownAction", []string{"JWT", "ID", "Type", "Payload"}, handler)
	assert.NotNil(t, form)

	// Check that form has expected fields
	jwtField := form.GetFormItemByLabel("JWT")
	assert.NotNil(t, jwtField)

	idField := form.GetFormItemByLabel("ID")
	assert.NotNil(t, idField)

	typeField := form.GetFormItemByLabel("Тип")
	assert.NotNil(t, typeField)

	payloadField := form.GetFormItemByLabel("Payload")
	assert.NotNil(t, payloadField)
}

// TestCreateDataFormWithUnknownField tests createDataForm with unknown field.
func TestCreateDataFormWithUnknownField(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	handler := func(values map[string]string) error {
		return nil
	}

	form := app.createDataForm("UnknownAction", []string{"UnknownField"}, handler)
	assert.NotNil(t, form)
}

// TestShowModalWithEmptyMessage tests showModal with empty message.
func TestShowModalWithEmptyMessage(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Test that showModal doesn't panic with empty message
	assert.NotPanics(t, func() {
		app.showModal("")
	})
}

// TestShowModalWithSpecialCharacters tests showModal with special characters.
func TestShowModalWithSpecialCharacters(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Test that showModal doesn't panic with special characters
	assert.NotPanics(t, func() {
		app.showModal("Test message with special chars: !@#$%^&*()")
	})
}

// TestTUIAppWithNilClients tests TUIApp behavior with nil clients.
func TestTUIAppWithNilClients(t *testing.T) {
	app := &TUIApp{
		app:        tview.NewApplication(),
		userClient: nil,
		dataClient: nil,
	}

	// Test that methods don't panic with nil clients
	assert.NotPanics(t, func() {
		app.showModal("Test message")
	})

	// Test that createAuthForm doesn't panic
	handler := func(data *AuthFormData) error {
		return nil
	}
	form := app.createAuthForm("Test", handler)
	assert.NotNil(t, form)

	// Test that createDataForm doesn't panic
	dataHandler := func(values map[string]string) error {
		return nil
	}
	dataForm := app.createDataForm("Test", []string{"JWT"}, dataHandler)
	assert.NotNil(t, dataForm)
}

// TestHandleRegisterWithNilClient tests registration with nil client.
func TestHandleRegisterWithNilClient(t *testing.T) {
	app := &TUIApp{
		userClient: nil,
		app:        tview.NewApplication(),
	}

	err := app.handleRegister(createTestAuthFormData())
	assert.Error(t, err)
}

// TestHandleLoginWithNilClient tests login with nil client.
func TestHandleLoginWithNilClient(t *testing.T) {
	app := &TUIApp{
		userClient: nil,
		app:        tview.NewApplication(),
	}

	err := app.handleLogin(createTestAuthFormData())
	assert.Error(t, err)
}

// TestHandleAddDataWithNilClient tests adding data with nil client.
func TestHandleAddDataWithNilClient(t *testing.T) {
	app := &TUIApp{
		dataClient: nil,
		app:        tview.NewApplication(),
	}

	err := app.handleAddData(createTestAddDataFormData())
	assert.Error(t, err)
}

// TestHandleGetDataWithNilClient tests getting data with nil client.
func TestHandleGetDataWithNilClient(t *testing.T) {
	app := &TUIApp{
		dataClient: nil,
		app:        tview.NewApplication(),
	}

	err := app.handleGetData(createTestJWTIDFormData())
	assert.Error(t, err)
}

// TestHandleEditDataWithNilClient tests editing data with nil client.
func TestHandleEditDataWithNilClient(t *testing.T) {
	app := &TUIApp{
		dataClient: nil,
		app:        tview.NewApplication(),
	}

	err := app.handleEditData(createTestEditDataFormData())
	assert.Error(t, err)
}

// TestHandleDeleteDataWithNilClient tests deleting data with nil client.
func TestHandleDeleteDataWithNilClient(t *testing.T) {
	app := &TUIApp{
		dataClient: nil,
		app:        tview.NewApplication(),
	}

	err := app.handleDeleteData(createTestJWTIDFormData())
	assert.Error(t, err)
}

// TestCreateDataFormWithAllFieldTypes tests createDataForm with all field types.
func TestCreateDataFormWithAllFieldTypes(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	handler := func(values map[string]string) error {
		return nil
	}

	// Test with all possible field types
	fields := []string{"JWT", "ID", "Type", "Payload"}
	form := app.createDataForm("TestAction", fields, handler)
	assert.NotNil(t, form)

	// Check that all fields are present (note: Type field has Russian label "Тип")
	assert.NotNil(t, form.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form.GetFormItemByLabel("ID"))
	assert.NotNil(t, form.GetFormItemByLabel("Тип")) // Russian label for Type
	assert.NotNil(t, form.GetFormItemByLabel("Payload"))
}

// TestCreateDataFormWithSpecificFields tests createDataForm with specific field combinations.
func TestCreateDataFormWithSpecificFields(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	handler := func(values map[string]string) error {
		return nil
	}

	// Test with JWT and ID fields (like get/delete forms)
	form1 := app.createDataForm("GetAction", []string{"JWT", "ID"}, handler)
	assert.NotNil(t, form1)
	assert.NotNil(t, form1.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form1.GetFormItemByLabel("ID"))

	// Test with JWT, Type, and Payload fields (like add form)
	form2 := app.createDataForm("AddAction", []string{"JWT", "Type", "Payload"}, handler)
	assert.NotNil(t, form2)
	assert.NotNil(t, form2.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form2.GetFormItemByLabel("Тип")) // Russian label for Type
	assert.NotNil(t, form2.GetFormItemByLabel("Payload"))
}

// TestCreateDataFormWithEmptyFields tests createDataForm with empty fields slice.
func TestCreateDataFormWithEmptyFields(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	handler := func(values map[string]string) error {
		return nil
	}

	// Test with empty fields slice
	form := app.createDataForm("TestAction", []string{}, handler)
	assert.NotNil(t, form)
}

// TestCreateAuthFormWithErrorHandler tests createAuthForm with error handler.
func TestCreateAuthFormWithErrorHandler(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	errorHandler := func(data *AuthFormData) error {
		return assert.AnError
	}

	// Test register form with error handler
	registerForm := app.createAuthForm("Зарегистрироваться", errorHandler)
	assert.NotNil(t, registerForm)

	// Test login form with error handler
	loginForm := app.createAuthForm("Войти", errorHandler)
	assert.NotNil(t, loginForm)
}

// TestCreateDataFormWithErrorHandler tests createDataForm with error handler.
func TestCreateDataFormWithErrorHandler(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	errorHandler := func(values map[string]string) error {
		return assert.AnError
	}

	// Test with error handler
	form := app.createDataForm("TestAction", []string{"JWT", "ID"}, errorHandler)
	assert.NotNil(t, form)
}

// TestShowModalWithLongMessage tests showModal with long message.
func TestShowModalWithLongMessage(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Test that showModal doesn't panic with long message
	longMessage := "This is a very long message that should test the modal's ability to handle long text without causing any issues or panics. It contains multiple sentences and should be displayed properly in the modal dialog."
	assert.NotPanics(t, func() {
		app.showModal(longMessage)
	})
}

// TestShowModalWithUnicodeCharacters tests showModal with unicode characters.
func TestShowModalWithUnicodeCharacters(t *testing.T) {
	app := &TUIApp{
		app: tview.NewApplication(),
	}

	// Test that showModal doesn't panic with unicode characters
	unicodeMessage := "Тест с русскими символами: привет мир! 🌍 测试中文 🎉"
	assert.NotPanics(t, func() {
		app.showModal(unicodeMessage)
	})
}

// TestRunTUI tests the RunTUI function.
func TestRunTUI(t *testing.T) {
	// Skip this test if we're in a CI environment or if tview can't initialize
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping RunTUI test in CI environment")
	}

	// Test that RunTUI doesn't panic immediately
	assert.NotPanics(t, func() {
		// This will likely fail due to missing gRPC connection, but shouldn't panic
		_ = RunTUI()
	})
}
