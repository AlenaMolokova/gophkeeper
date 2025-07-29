// Package tui provides a Terminal User Interface (TUI) for the GophKeeper client.
// This file contains form data extractors to convert tview.Form input fields
// to typed structures, ensuring type safety and validation.
package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

// extractInputFieldText safely extracts text from a form input field by label.
// It returns an error if the field is not found or is not an InputField.
func extractInputFieldText(form *tview.Form, label string) (string, error) {
	item := form.GetFormItemByLabel(label)
	if item == nil {
		return "", fmt.Errorf("form field '%s' not found", label)
	}

	inputField, ok := item.(*tview.InputField)
	if !ok {
		return "", fmt.Errorf("form field '%s' is not an input field", label)
	}

	return inputField.GetText(), nil
}

// extractPasswordFieldText safely extracts text from a form password field by label.
// It returns an error if the field is not found or is not a PasswordField.
func extractPasswordFieldText(form *tview.Form, label string) (string, error) {
	item := form.GetFormItemByLabel(label)
	if item == nil {
		return "", fmt.Errorf("form field '%s' not found", label)
	}

	passwordField, ok := item.(*tview.InputField) // PasswordField is also InputField
	if !ok {
		return "", fmt.Errorf("form field '%s' is not a password field", label)
	}

	return passwordField.GetText(), nil
}

// ExtractAuthFormData extracts authentication form data from a tview.Form.
// It validates the extracted data and returns a typed AuthFormData structure.
func ExtractAuthFormData(form *tview.Form) (*AuthFormData, error) {
	email, err := extractInputFieldText(form, "Email")
	if err != nil {
		return nil, fmt.Errorf("failed to extract email: %w", err)
	}

	password, err := extractPasswordFieldText(form, "Пароль")
	if err != nil {
		return nil, fmt.Errorf("failed to extract password: %w", err)
	}

	data := &AuthFormData{
		Email:    email,
		Password: password,
	}

	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return data, nil
}

// ExtractJWTIDFormData extracts JWT ID form data from a tview.Form.
// It validates the extracted data and returns a typed JWTIDFormData structure.
func ExtractJWTIDFormData(form *tview.Form) (*JWTIDFormData, error) {
	jwt, err := extractInputFieldText(form, "JWT")
	if err != nil {
		return nil, fmt.Errorf("failed to extract JWT: %w", err)
	}

	id, err := extractInputFieldText(form, "ID")
	if err != nil {
		return nil, fmt.Errorf("failed to extract ID: %w", err)
	}

	data := &JWTIDFormData{
		JWT: jwt,
		ID:  id,
	}

	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return data, nil
}

// ExtractAddDataFormData extracts add data form data from a tview.Form.
// It validates the extracted data and returns a typed AddDataFormData structure.
func ExtractAddDataFormData(form *tview.Form) (*AddDataFormData, error) {
	jwt, err := extractInputFieldText(form, "JWT")
	if err != nil {
		return nil, fmt.Errorf("failed to extract JWT: %w", err)
	}

	dataType, err := extractInputFieldText(form, "Тип")
	if err != nil {
		return nil, fmt.Errorf("failed to extract type: %w", err)
	}

	payload, err := extractInputFieldText(form, "Payload")
	if err != nil {
		return nil, fmt.Errorf("failed to extract payload: %w", err)
	}

	data := &AddDataFormData{
		JWT:     jwt,
		Type:    dataType,
		Payload: payload,
	}

	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return data, nil
}

// ExtractEditDataFormData extracts edit data form data from a tview.Form.
// It validates the extracted data and returns a typed EditDataFormData structure.
func ExtractEditDataFormData(form *tview.Form) (*EditDataFormData, error) {
	jwt, err := extractInputFieldText(form, "JWT")
	if err != nil {
		return nil, fmt.Errorf("failed to extract JWT: %w", err)
	}

	id, err := extractInputFieldText(form, "ID")
	if err != nil {
		return nil, fmt.Errorf("failed to extract ID: %w", err)
	}

	dataType, err := extractInputFieldText(form, "Тип")
	if err != nil {
		return nil, fmt.Errorf("failed to extract type: %w", err)
	}

	payload, err := extractInputFieldText(form, "Payload")
	if err != nil {
		return nil, fmt.Errorf("failed to extract payload: %w", err)
	}

	data := &EditDataFormData{
		JWT:     jwt,
		ID:      id,
		Type:    dataType,
		Payload: payload,
	}

	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return data, nil
}
