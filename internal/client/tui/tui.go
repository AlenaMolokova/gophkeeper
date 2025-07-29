// Package tui provides a Terminal User Interface (TUI) for the GophKeeper client.
// It implements an interactive command-line interface using the tview library
// for user authentication, data management, and server communication.
//
// The package provides a complete user interface with forms, menus, and modals
// for managing encrypted data through a secure gRPC connection to the server.
package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	clientapplogic "github.com/AlenaMolokova/gophkeeper/internal/client/applogic"
	"github.com/AlenaMolokova/gophkeeper/internal/config"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// TUIApp represents the TUI application with its components.
// It manages the main application window, menu system, and gRPC client connections
// for both user authentication and data management services.
type TUIApp struct {
	app        *tview.Application          // Main application instance
	menu       *tview.List                 // Main menu component
	userClient clientapi.UserServiceClient // Client for user authentication
	dataClient clientapi.DataServiceClient // Client for data operations
	conn       *grpc.ClientConn            // gRPC connection to server
	localizer  *Localizer                  // Localization manager
}

// showModal displays a modal dialog with the given message.
// It creates a modal window with an "OK" button and returns to the main menu
// when the user dismisses the dialog.
func (t *TUIApp) showModal(message string) {
	modal := tview.NewModal().SetText(message).AddButtons([]string{"OK"}).SetDoneFunc(func(_ int, _ string) {
		t.app.SetRoot(t.menu, true)
	})
	t.app.SetRoot(modal, true)
}

// handleError displays an error message in a modal dialog.
// This function eliminates code duplication by providing a centralized
// way to handle and display errors to the user.
func (t *TUIApp) handleError(err error) {
	t.showModal(fmt.Sprintf("Ошибка: %v", err))
}

// createTimeoutContext creates a context with a 5-second timeout.
// This function eliminates code duplication by providing a centralized
// way to create timeout contexts for gRPC operations.
func (t *TUIApp) createTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// mapJWTIDFormData converts JWTIDFormData to map[string]string for backward compatibility.
// This function eliminates code duplication in form data mapping.
func mapJWTIDFormData(data *JWTIDFormData) map[string]string {
	return map[string]string{
		"JWT": data.JWT,
		"ID":  data.ID,
	}
}

// mapAddDataFormData converts AddDataFormData to map[string]string for backward compatibility.
// This function eliminates code duplication in form data mapping.
func mapAddDataFormData(data *AddDataFormData) map[string]string {
	return map[string]string{
		"JWT":     data.JWT,
		"Type":    data.Type,
		"Payload": data.Payload,
	}
}

// mapEditDataFormData converts EditDataFormData to map[string]string for backward compatibility.
// This function eliminates code duplication in form data mapping.
func mapEditDataFormData(data *EditDataFormData) map[string]string {
	return map[string]string{
		"JWT":     data.JWT,
		"ID":      data.ID,
		"Type":    data.Type,
		"Payload": data.Payload,
	}
}

// createAuthForm creates a form for registration or login.
// It delegates to the appropriate form constructor based on the action.
func (t *TUIApp) createAuthForm(action string, handler FormHandler[*AuthFormData]) *tview.Form {
	errorHandler := func(data *AuthFormData) error {
		if err := handler(data); err != nil {
			t.handleError(err)
			return err
		}
		return nil
	}

	backHandler := func() { t.app.SetRoot(t.menu, true) }

	switch action {
	case "Зарегистрироваться":
		return createRegisterForm(errorHandler, backHandler)
	case "Войти":
		return createLoginForm(errorHandler, backHandler)
	default:
		// Fallback to generic form
		form := tview.NewForm()
		f := form
		form.
			AddInputField("Email", "", 30, nil, nil).
			AddPasswordField("Пароль", "", 30, '*', nil).
			AddButton(action, func() {
				data, err := ExtractAuthFormData(f)
				if err != nil {
					t.handleError(err)
					return
				}
				_ = errorHandler(data)
			}).
			AddButton("Назад", backHandler)
		return form
	}
}

// handleRegister handles user registration.
// It creates a context with timeout, calls the registration logic,
// and displays the result to the user via a modal dialog.
func (t *TUIApp) handleRegister(data *AuthFormData) error {
	ctx, cancel := t.createTimeoutContext()
	defer cancel()

	token, err := clientapplogic.RegisterUser(ctx, t.userClient, data.Email, data.Password)
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Успешно! JWT: %s", token))
	return nil
}

// handleLogin handles user login.
// It creates a context with timeout, calls the login logic,
// and displays the result to the user via a modal dialog.
func (t *TUIApp) handleLogin(data *AuthFormData) error {
	ctx, cancel := t.createTimeoutContext()
	defer cancel()

	token, err := clientapplogic.LoginUser(ctx, t.userClient, data.Email, data.Password)
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Успешно! JWT: %s", token))
	return nil
}

// createDataForm creates a form for data operations.
// createDataForm creates a form for data operations based on the action type.
// It delegates to the appropriate form constructor based on the action and fields.
func (t *TUIApp) createDataForm(action string, fields []string, handler func(map[string]string) error) *tview.Form {
	backHandler := func() { t.app.SetRoot(t.menu, true) }

	// Use action-specific form creators
	switch action {
	case "Получить":
		return t.createGetDataForm(handler, backHandler)
	case "Добавить":
		return t.createAddDataForm(handler, backHandler)
	case "Изменить":
		return t.createEditDataForm(handler, backHandler)
	case "Удалить":
		return t.createDeleteDataForm(handler, backHandler)
	default:
		return t.createGenericForm(action, fields, handler, backHandler)
	}
}

// createGetDataForm creates a form for getting data.
func (t *TUIApp) createGetDataForm(handler func(map[string]string) error, backHandler func()) *tview.Form {
	errorHandler := func(data *JWTIDFormData) error {
		values := mapJWTIDFormData(data)
		if err := handler(values); err != nil {
			t.handleError(err)
			return err
		}
		return nil
	}
	return createGetDataForm(errorHandler, backHandler)
}

// createAddDataForm creates a form for adding data.
func (t *TUIApp) createAddDataForm(handler func(map[string]string) error, backHandler func()) *tview.Form {
	errorHandler := func(data *AddDataFormData) error {
		values := mapAddDataFormData(data)
		if err := handler(values); err != nil {
			t.handleError(err)
			return err
		}
		return nil
	}
	return createAddDataForm(errorHandler, backHandler)
}

// createEditDataForm creates a form for editing data.
func (t *TUIApp) createEditDataForm(handler func(map[string]string) error, backHandler func()) *tview.Form {
	errorHandler := func(data *EditDataFormData) error {
		values := mapEditDataFormData(data)
		if err := handler(values); err != nil {
			t.handleError(err)
			return err
		}
		return nil
	}
	return createEditDataForm(errorHandler, backHandler)
}

// createDeleteDataForm creates a form for deleting data.
func (t *TUIApp) createDeleteDataForm(handler func(map[string]string) error, backHandler func()) *tview.Form {
	errorHandler := func(data *JWTIDFormData) error {
		values := mapJWTIDFormData(data)
		if err := handler(values); err != nil {
			t.handleError(err)
			return err
		}
		return nil
	}
	return createDeleteDataForm(errorHandler, backHandler)
}

// createGenericForm creates a generic form for unknown actions.
func (t *TUIApp) createGenericForm(action string, fields []string, handler func(map[string]string) error, backHandler func()) *tview.Form {
	form := tview.NewForm()
	f := form

	// Add form fields based on the field type
	t.addFormFields(form, fields)

	// Add action and back buttons
	form.
		AddButton(action, func() {
			values := t.extractFormValues(f, fields)
			if err := handler(values); err != nil {
				t.handleError(err)
			}
		}).
		AddButton("Назад", backHandler)
	return form
}

// addFormFields adds input fields to the form based on field types.
func (t *TUIApp) addFormFields(form *tview.Form, fields []string) {
	for _, field := range fields {
		switch field {
		case "JWT":
			form.AddInputField("JWT", "", 60, nil, nil)
		case "ID":
			form.AddInputField("ID", "", 36, nil, nil)
		case "Type":
			form.AddInputField("Тип", "", 10, nil, nil)
		case "Payload":
			form.AddInputField("Payload", "", 60, nil, nil)
		}
	}
}

// extractFormValues extracts values from form fields.
func (t *TUIApp) extractFormValues(form *tview.Form, fields []string) map[string]string {
	values := make(map[string]string)
	for _, field := range fields {
		item := form.GetFormItemByLabel(field)
		if input, ok := item.(*tview.InputField); ok {
			values[field] = input.GetText()
		}
	}
	return values
}

// handleGetData handles retrieving data by ID.
// It creates a context with timeout, calls the data retrieval logic,
// and displays the retrieved data to the user via a modal dialog.
func (t *TUIApp) handleGetData(data *JWTIDFormData) error {
	ctx, cancel := t.createTimeoutContext()
	defer cancel()

	result, payload, err := clientapplogic.GetData(ctx, t.dataClient, data.JWT, data.ID)
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Данные получены:\nID: %s\nТип: %s\nPayload: %s", result.Id, result.Type, payload))
	return nil
}

// handleAddData handles adding new data.
// It creates a context with timeout, calls the data addition logic,
// and displays the result to the user via a modal dialog.
func (t *TUIApp) handleAddData(data *AddDataFormData) error {
	ctx, cancel := t.createTimeoutContext()
	defer cancel()

	id, err := clientapplogic.AddData(ctx, t.dataClient, data.JWT, data.Type, data.Payload)
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Данные добавлены. ID: %s", id))
	return nil
}

// handleEditData handles editing existing data.
// It creates a context with timeout, calls the data editing logic,
// and displays the result to the user via a modal dialog.
func (t *TUIApp) handleEditData(data *EditDataFormData) error {
	ctx, cancel := t.createTimeoutContext()
	defer cancel()

	newID, err := clientapplogic.EditData(ctx, t.dataClient, data.JWT, data.ID, data.Type, data.Payload)
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Данные обновлены. ID: %s", newID))
	return nil
}

// handleDeleteData handles deleting data by ID.
// It creates a context with timeout, calls the data deletion logic,
// and displays the result to the user via a modal dialog.
func (t *TUIApp) handleDeleteData(data *JWTIDFormData) error {
	ctx, cancel := t.createTimeoutContext()
	defer cancel()

	err := clientapplogic.DeleteData(ctx, t.dataClient, data.JWT, data.ID)
	if err != nil {
		return err
	}

	t.showModal("Данные удалены.")
	return nil
}

// setupMenu creates the main menu with all available actions.
// It initializes the main menu with navigation items for all available
// operations including user authentication and data management.
func (t *TUIApp) setupMenu() {
	t.menu = tview.NewList().
		AddItem(t.localizer.Translate("registration"), t.localizer.Translate("registration_description"), 'r', func() {
			form := t.createAuthForm(t.localizer.Translate("register"), t.handleRegister)
			t.app.SetRoot(form, true)
		}).
		AddItem(t.localizer.Translate("login"), t.localizer.Translate("login_description"), 'l', func() {
			form := t.createAuthForm(t.localizer.Translate("login_action"), t.handleLogin)
			t.app.SetRoot(form, true)
		}).
		AddItem(t.localizer.Translate("view_data"), t.localizer.Translate("view_data_description"), 'g', func() {
			// Adapter for backward compatibility
			adapter := func(values map[string]string) error {
				data := &JWTIDFormData{
					JWT: values["JWT"],
					ID:  values["ID"],
				}
				return t.handleGetData(data)
			}
			form := t.createDataForm(t.localizer.Translate("get"), []string{"JWT", "ID"}, adapter)
			t.app.SetRoot(form, true)
		}).
		AddItem(t.localizer.Translate("add_data"), t.localizer.Translate("add_data_description"), 'a', func() {
			// Adapter for backward compatibility
			adapter := func(values map[string]string) error {
				data := &AddDataFormData{
					JWT:     values["JWT"],
					Type:    values["Type"],
					Payload: values["Payload"],
				}
				return t.handleAddData(data)
			}
			form := t.createDataForm(t.localizer.Translate("add"), []string{"JWT", "Type", "Payload"}, adapter)
			t.app.SetRoot(form, true)
		}).
		AddItem(t.localizer.Translate("edit_data"), t.localizer.Translate("edit_data_description"), 'e', func() {
			// Adapter for backward compatibility
			adapter := func(values map[string]string) error {
				data := &EditDataFormData{
					JWT:     values["JWT"],
					ID:      values["ID"],
					Type:    values["Type"],
					Payload: values["Payload"],
				}
				return t.handleEditData(data)
			}
			form := t.createDataForm(t.localizer.Translate("edit"), []string{"JWT", "ID", "Type", "Payload"}, adapter)
			t.app.SetRoot(form, true)
		}).
		AddItem(t.localizer.Translate("delete_data"), t.localizer.Translate("delete_data_description"), 'd', func() {
			// Adapter for backward compatibility
			adapter := func(values map[string]string) error {
				data := &JWTIDFormData{
					JWT: values["JWT"],
					ID:  values["ID"],
				}
				return t.handleDeleteData(data)
			}
			form := t.createDataForm(t.localizer.Translate("delete"), []string{"JWT", "ID"}, adapter)
			t.app.SetRoot(form, true)
		}).
		AddItem(t.localizer.Translate("exit"), t.localizer.Translate("exit_description"), 'q', func() { t.app.Stop() })
}

// RunTUI starts the TUI application.
// It initializes the gRPC connection to the server, creates the TUI application,
// sets up the main menu, and starts the interactive interface.
//
// The function handles the complete application lifecycle including
// connection setup, UI initialization, and graceful shutdown.
func RunTUI() error {
	config := config.NewClientConfig()

	// Initialize localization
	localizer, err := NewLocalizer()
	if err != nil {
		return err
	}

	creds, err := credentials.NewClientTLSFromFile(config.CertPath, "")
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(config.GetClientServerAddress(), grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()

	userClient := clientapi.NewUserServiceClient(conn)
	dataClient := clientapi.NewDataServiceClient(conn)

	tui := &TUIApp{
		app:        tview.NewApplication(),
		userClient: userClient,
		dataClient: dataClient,
		conn:       conn,
		localizer:  localizer,
	}

	tui.setupMenu()
	return tui.app.SetRoot(tui.menu, true).Run()
}
