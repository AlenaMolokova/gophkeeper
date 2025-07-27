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

// createAuthForm creates a form for registration or login.
// It generates a form with email and password fields, along with action buttons
// for submitting the form or returning to the main menu.
//
// The handler function is called with the form data when the user submits.
func (t *TUIApp) createAuthForm(action string, handler func(email, password string) error) *tview.Form {
	form := tview.NewForm()
	f := form // save reference for closure

	form.
		AddInputField("Email", "", 30, nil, nil).
		AddPasswordField("Пароль", "", 30, '*', nil).
		AddButton(action, func() {
			email := f.GetFormItemByLabel("Email").(*tview.InputField).GetText()
			password := f.GetFormItemByLabel("Пароль").(*tview.InputField).GetText()
			if err := handler(email, password); err != nil {
				t.showModal(fmt.Sprintf("Ошибка: %v", err))
				return
			}
		}).
		AddButton("Назад", func() { t.app.SetRoot(t.menu, true) })

	return form
}

// handleRegister handles user registration.
// It creates a context with timeout, calls the registration logic,
// and displays the result to the user via a modal dialog.
func (t *TUIApp) handleRegister(email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := clientapplogic.RegisterUser(ctx, t.userClient, email, password)
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Успешно! JWT: %s", token))
	return nil
}

// handleLogin handles user login.
// It creates a context with timeout, calls the login logic,
// and displays the result to the user via a modal dialog.
func (t *TUIApp) handleLogin(email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := clientapplogic.LoginUser(ctx, t.userClient, email, password)
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Успешно! JWT: %s", token))
	return nil
}

// createDataForm creates a form for data operations.
// It generates a dynamic form based on the required fields for the operation,
// including JWT token, data ID, type, and payload fields as needed.
//
// The handler function is called with the form data when the user submits.
func (t *TUIApp) createDataForm(action string, fields []string, handler func(map[string]string) error) *tview.Form {
	form := tview.NewForm()
	f := form // save reference for closure

	// Add input fields based on the action
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

	form.
		AddButton(action, func() {
			values := make(map[string]string)
			for _, field := range fields {
				item := f.GetFormItemByLabel(field)
				if input, ok := item.(*tview.InputField); ok {
					values[field] = input.GetText()
				}
			}
			if err := handler(values); err != nil {
				t.showModal(fmt.Sprintf("Ошибка: %v", err))
				return
			}
		}).
		AddButton("Назад", func() { t.app.SetRoot(t.menu, true) })

	return form
}

// handleGetData handles retrieving data by ID.
// It creates a context with timeout, calls the data retrieval logic,
// and displays the retrieved data to the user via a modal dialog.
func (t *TUIApp) handleGetData(values map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, payload, err := clientapplogic.GetData(ctx, t.dataClient, values["JWT"], values["ID"])
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Данные получены:\nID: %s\nТип: %s\nPayload: %s", data.Id, data.Type, payload))
	return nil
}

// handleAddData handles adding new data.
// It creates a context with timeout, calls the data addition logic,
// and displays the result to the user via a modal dialog.
func (t *TUIApp) handleAddData(values map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := clientapplogic.AddData(ctx, t.dataClient, values["JWT"], values["Type"], values["Payload"])
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Данные добавлены. ID: %s", id))
	return nil
}

// handleEditData handles editing existing data.
// It creates a context with timeout, calls the data editing logic,
// and displays the result to the user via a modal dialog.
func (t *TUIApp) handleEditData(values map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newID, err := clientapplogic.EditData(ctx, t.dataClient, values["JWT"], values["ID"], values["Type"], values["Payload"])
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Данные обновлены. ID: %s", newID))
	return nil
}

// handleDeleteData handles deleting data by ID.
// It creates a context with timeout, calls the data deletion logic,
// and displays the result to the user via a modal dialog.
func (t *TUIApp) handleDeleteData(values map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := clientapplogic.DeleteData(ctx, t.dataClient, values["JWT"], values["ID"]); err != nil {
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
		AddItem("Регистрация", "Создать нового пользователя", 'r', func() {
			form := t.createAuthForm("Зарегистрироваться", t.handleRegister)
			t.app.SetRoot(form, true)
		}).
		AddItem("Вход", "Войти в систему", 'l', func() {
			form := t.createAuthForm("Войти", t.handleLogin)
			t.app.SetRoot(form, true)
		}).
		AddItem("Просмотр данных", "Получить данные по ID", 'g', func() {
			form := t.createDataForm("Получить", []string{"JWT", "ID"}, t.handleGetData)
			t.app.SetRoot(form, true)
		}).
		AddItem("Добавить данные", "Добавить новую запись", 'a', func() {
			form := t.createDataForm("Добавить", []string{"JWT", "Type", "Payload"}, t.handleAddData)
			t.app.SetRoot(form, true)
		}).
		AddItem("Редактировать данные", "Изменить существующую запись", 'e', func() {
			form := t.createDataForm("Изменить", []string{"JWT", "ID", "Type", "Payload"}, t.handleEditData)
			t.app.SetRoot(form, true)
		}).
		AddItem("Удалить данные", "Удалить запись по ID", 'd', func() {
			form := t.createDataForm("Удалить", []string{"JWT", "ID"}, t.handleDeleteData)
			t.app.SetRoot(form, true)
		}).
		AddItem("Выход", "Завершить работу", 'q', func() { t.app.Stop() })
}

// RunTUI starts the TUI application.
// It initializes the gRPC connection to the server, creates the TUI application,
// sets up the main menu, and starts the interactive interface.
//
// The function handles the complete application lifecycle including
// connection setup, UI initialization, and graceful shutdown.
func RunTUI() error {
	config := config.NewClientConfig()

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
	}

	tui.setupMenu()
	return tui.app.SetRoot(tui.menu, true).Run()
}
