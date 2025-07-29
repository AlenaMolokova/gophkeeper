// Package tui provides a Terminal User Interface (TUI) for the GophKeeper client.
// It implements an interactive command-line interface using the tview library
// for user authentication, data management, and server communication.
//
// The package provides a complete user interface with forms, menus, and modals
// for managing encrypted data through a secure gRPC connection to the server.
package tui

import (
	"github.com/rivo/tview"
)

// createAuthFormBase creates a base form for authentication operations.
// It generates a form with email and password fields and customizable action button.
// The form uses typed AuthFormData for type safety and validation.
func createAuthFormBase(action string, handler FormHandler[*AuthFormData], backHandler func()) *tview.Form {
	form := tview.NewForm()
	f := form // save reference for closure

	form.
		AddInputField("Email", "", 30, nil, nil).
		AddPasswordField("Пароль", "", 30, '*', nil).
		AddButton(action, func() {
			data, err := ExtractAuthFormData(f)
			if err != nil {
				// Error handling is done by the caller
				return
			}
			if err := handler(data); err != nil {
				// Error handling is done by the caller
				return
			}
		}).
		AddButton("Назад", backHandler)

	return form
}

// createLoginForm creates a form for user login.
// It generates a form with email and password fields for user authentication.
func createLoginForm(handler FormHandler[*AuthFormData], backHandler func()) *tview.Form {
	return createAuthFormBase("Войти", handler, backHandler)
}

// createRegisterForm creates a form for user registration.
// It generates a form with email and password fields for user registration.
func createRegisterForm(handler FormHandler[*AuthFormData], backHandler func()) *tview.Form {
	return createAuthFormBase("Зарегистрироваться", handler, backHandler)
}

// createJWTIDFormBase creates a base form for operations requiring JWT and ID.
// It generates a form with JWT token and data ID fields and customizable action button.
// The form uses typed JWTIDFormData for type safety and validation.
func createJWTIDFormBase(action string, handler FormHandler[*JWTIDFormData], backHandler func()) *tview.Form {
	form := tview.NewForm()
	f := form // save reference for closure

	form.
		AddInputField("JWT", "", 60, nil, nil).
		AddInputField("ID", "", 36, nil, nil).
		AddButton(action, func() {
			data, err := ExtractJWTIDFormData(f)
			if err != nil {
				// Error handling is done by the caller
				return
			}
			if err := handler(data); err != nil {
				// Error handling is done by the caller
				return
			}
		}).
		AddButton("Назад", backHandler)

	return form
}

// createGetDataForm creates a form for retrieving data by ID.
// It generates a form with JWT token and data ID fields.
func createGetDataForm(handler FormHandler[*JWTIDFormData], backHandler func()) *tview.Form {
	return createJWTIDFormBase("Получить", handler, backHandler)
}

// createAddDataForm creates a form for adding new data.
// It generates a form with JWT token, data type, and payload fields.
// The form uses typed AddDataFormData for type safety and validation.
func createAddDataForm(handler FormHandler[*AddDataFormData], backHandler func()) *tview.Form {
	form := tview.NewForm()
	f := form // save reference for closure

	form.
		AddInputField("JWT", "", 60, nil, nil).
		AddInputField("Тип", "", 10, nil, nil).
		AddInputField("Payload", "", 60, nil, nil).
		AddButton("Добавить", func() {
			data, err := ExtractAddDataFormData(f)
			if err != nil {
				// Error handling is done by the caller
				return
			}
			if err := handler(data); err != nil {
				// Error handling is done by the caller
				return
			}
		}).
		AddButton("Назад", backHandler)

	return form
}

// createEditDataForm creates a form for editing existing data.
// It generates a form with JWT token, data ID, type, and payload fields.
// The form uses typed EditDataFormData for type safety and validation.
func createEditDataForm(handler FormHandler[*EditDataFormData], backHandler func()) *tview.Form {
	form := tview.NewForm()
	f := form // save reference for closure

	form.
		AddInputField("JWT", "", 60, nil, nil).
		AddInputField("ID", "", 36, nil, nil).
		AddInputField("Тип", "", 10, nil, nil).
		AddInputField("Payload", "", 60, nil, nil).
		AddButton("Изменить", func() {
			data, err := ExtractEditDataFormData(f)
			if err != nil {
				// Error handling is done by the caller
				return
			}
			if err := handler(data); err != nil {
				// Error handling is done by the caller
				return
			}
		}).
		AddButton("Назад", backHandler)

	return form
}

// createDeleteDataForm creates a form for deleting data by ID.
// It generates a form with JWT token and data ID fields.
func createDeleteDataForm(handler FormHandler[*JWTIDFormData], backHandler func()) *tview.Form {
	return createJWTIDFormBase("Удалить", handler, backHandler)
}
