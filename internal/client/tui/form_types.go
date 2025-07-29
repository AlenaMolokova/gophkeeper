// Package tui provides a Terminal User Interface (TUI) for the GophKeeper client.
// This file contains typed structures for form data to replace map[string]string
// and ensure type safety according to pravila.md principles.
package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

// FormData represents the base interface for all form data structures.
// It provides a common contract for form data validation and processing.
type FormData interface {
	// Validate checks if the form data is valid and returns an error if not.
	Validate() error
}

// AuthFormData represents authentication form data (login/register).
// It contains email and password fields for user authentication.
type AuthFormData struct {
	Email    string // User's email address
	Password string // User's password
}

// Validate checks if the authentication form data is valid.
// It ensures that both email and password are provided.
func (a *AuthFormData) Validate() error {
	if a.Email == "" {
		return fmt.Errorf("email is required")
	}
	if a.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

// JWTIDFormData represents form data for operations requiring JWT token and data ID.
// It is used for get and delete operations that need both authentication and data identification.
type JWTIDFormData struct {
	JWT string // JWT authentication token
	ID  string // Data identifier
}

// Validate checks if the JWT ID form data is valid.
// It ensures that both JWT token and ID are provided.
func (j *JWTIDFormData) Validate() error {
	if j.JWT == "" {
		return fmt.Errorf("JWT token is required")
	}
	if j.ID == "" {
		return fmt.Errorf("data ID is required")
	}
	return nil
}

// AddDataFormData represents form data for adding new data entries.
// It contains JWT token, data type, and payload for creating new data.
type AddDataFormData struct {
	JWT     string // JWT authentication token
	Type    string // Data type (login, text, binary, card, otp)
	Payload string // Data payload to be encrypted and stored
}

// Validate checks if the add data form data is valid.
// It ensures that all required fields are provided and data type is valid.
func (a *AddDataFormData) Validate() error {
	if a.JWT == "" {
		return fmt.Errorf("JWT token is required")
	}
	if a.Type == "" {
		return fmt.Errorf("data type is required")
	}
	if a.Payload == "" {
		return fmt.Errorf("payload is required")
	}

	// Validate data type
	validTypes := map[string]bool{
		"login":  true,
		"text":   true,
		"binary": true,
		"card":   true,
		"otp":    true,
	}
	if !validTypes[a.Type] {
		return fmt.Errorf("invalid data type: %s. Valid types are: login, text, binary, card, otp", a.Type)
	}

	return nil
}

// EditDataFormData represents form data for editing existing data entries.
// It contains JWT token, data ID, type, and payload for updating data.
type EditDataFormData struct {
	JWT     string // JWT authentication token
	ID      string // Data identifier
	Type    string // Data type (login, text, binary, card, otp)
	Payload string // Updated data payload
}

// Validate checks if the edit data form data is valid.
// It ensures that all required fields are provided and data type is valid.
func (e *EditDataFormData) Validate() error {
	if e.JWT == "" {
		return fmt.Errorf("JWT token is required")
	}
	if e.ID == "" {
		return fmt.Errorf("data ID is required")
	}
	if e.Type == "" {
		return fmt.Errorf("data type is required")
	}
	if e.Payload == "" {
		return fmt.Errorf("payload is required")
	}

	// Validate data type
	validTypes := map[string]bool{
		"login":  true,
		"text":   true,
		"binary": true,
		"card":   true,
		"otp":    true,
	}
	if !validTypes[e.Type] {
		return fmt.Errorf("invalid data type: %s. Valid types are: login, text, binary, card, otp", e.Type)
	}

	return nil
}

// FormHandler represents a generic form handler function.
// It provides type safety for form data processing.
type FormHandler[T FormData] func(data T) error

// FormDataExtractor represents a function that extracts form data from tview.Form.
// It provides a way to convert form input fields to typed structures.
type FormDataExtractor[T FormData] func(form *tview.Form) (T, error)
