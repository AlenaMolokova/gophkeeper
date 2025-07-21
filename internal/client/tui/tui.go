package tui

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	clientapplogic "github.com/AlenaMolokova/gophkeeper/internal/client/applogic"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// TUIApp represents the TUI application with its components.
type TUIApp struct {
	app    *tview.Application
	menu   *tview.List
	client clientapi.GophKeeperClient
	conn   *grpc.ClientConn
}

// showModal displays a modal dialog with the given message.
func (t *TUIApp) showModal(message string) {
	modal := tview.NewModal().SetText(message).AddButtons([]string{"OK"}).SetDoneFunc(func(_ int, _ string) {
		t.app.SetRoot(t.menu, true)
	})
	t.app.SetRoot(modal, true)
}

// createAuthForm creates a form for registration or login.
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
func (t *TUIApp) handleRegister(email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := clientapplogic.RegisterUser(ctx, t.client, email, password)
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Успешно! JWT: %s", token))
	return nil
}

// handleLogin handles user login.
func (t *TUIApp) handleLogin(email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := clientapplogic.LoginUser(ctx, t.client, email, password)
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Успешно! JWT: %s", token))
	return nil
}

// createDataForm creates a form for data operations.
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

	form.AddButton(action, func() {
		values := make(map[string]string)
		for _, field := range fields {
			switch field {
			case "JWT":
				values["JWT"] = f.GetFormItemByLabel("JWT").(*tview.InputField).GetText()
			case "ID":
				values["ID"] = f.GetFormItemByLabel("ID").(*tview.InputField).GetText()
			case "Type":
				values["Type"] = f.GetFormItemByLabel("Тип").(*tview.InputField).GetText()
			case "Payload":
				values["Payload"] = f.GetFormItemByLabel("Payload").(*tview.InputField).GetText()
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

// handleGetData handles getting data by ID.
func (t *TUIApp) handleGetData(values map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, decPayload, err := clientapplogic.GetData(ctx, t.client, values["JWT"], values["ID"])
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("ID: %s\nType: %s\nPayload: %s\nMetadata: %v\nTimestamp: %d",
		data.Id, data.Type, decPayload, data.Metadata, data.Timestamp))
	return nil
}

// handleAddData handles adding new data.
func (t *TUIApp) handleAddData(values map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := clientapplogic.AddData(ctx, t.client, values["JWT"], values["Type"], values["Payload"])
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Данные добавлены. ID: %s", id))
	return nil
}

// handleEditData handles editing existing data.
func (t *TUIApp) handleEditData(values map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newID, err := clientapplogic.EditData(ctx, t.client, values["JWT"], values["ID"], values["Type"], values["Payload"])
	if err != nil {
		return err
	}

	t.showModal(fmt.Sprintf("Данные обновлены. ID: %s", newID))
	return nil
}

// handleDeleteData handles deleting data by ID.
func (t *TUIApp) handleDeleteData(values map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := clientapplogic.DeleteData(ctx, t.client, values["JWT"], values["ID"]); err != nil {
		return err
	}

	t.showModal("Данные удалены.")
	return nil
}

// setupMenu creates the main menu with all available actions.
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
func RunTUI() error {
	certPath := os.Getenv("GOPHKEEPER_CERT")
	if certPath == "" {
		certPath = "cert/server.crt"
	}

	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := clientapi.NewGophKeeperClient(conn)

	tui := &TUIApp{
		app:    tview.NewApplication(),
		client: client,
		conn:   conn,
	}

	tui.setupMenu()
	return tui.app.SetRoot(tui.menu, true).Run()
}
