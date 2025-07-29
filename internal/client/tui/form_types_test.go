// Package tui provides a Terminal User Interface (TUI) for the GophKeeper client.
// This file contains tests for typed form data structures.
package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAuthFormData_Validate tests the validation of AuthFormData.
func TestAuthFormData_Validate(t *testing.T) {
	tests := []struct {
		name    string
		data    *AuthFormData
		wantErr bool
	}{
		{
			name: "valid data",
			data: &AuthFormData{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "empty email",
			data: &AuthFormData{
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "empty password",
			data: &AuthFormData{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			data: &AuthFormData{
				Email:    "",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestJWTIDFormData_Validate tests the validation of JWTIDFormData.
func TestJWTIDFormData_Validate(t *testing.T) {
	tests := []struct {
		name    string
		data    *JWTIDFormData
		wantErr bool
	}{
		{
			name: "valid data",
			data: &JWTIDFormData{
				JWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				ID:  "123e4567-e89b-12d3-a456-426614174000",
			},
			wantErr: false,
		},
		{
			name: "empty JWT",
			data: &JWTIDFormData{
				JWT: "",
				ID:  "123e4567-e89b-12d3-a456-426614174000",
			},
			wantErr: true,
		},
		{
			name: "empty ID",
			data: &JWTIDFormData{
				JWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				ID:  "",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			data: &JWTIDFormData{
				JWT: "",
				ID:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestAddDataFormData_Validate tests the validation of AddDataFormData.
func TestAddDataFormData_Validate(t *testing.T) {
	tests := []struct {
		name    string
		data    *AddDataFormData
		wantErr bool
	}{
		{
			name: "valid login data",
			data: &AddDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				Type:    "login",
				Payload: "username:test,password:secret",
			},
			wantErr: false,
		},
		{
			name: "valid text data",
			data: &AddDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				Type:    "text",
				Payload: "Some secret text",
			},
			wantErr: false,
		},
		{
			name: "valid binary data",
			data: &AddDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				Type:    "binary",
				Payload: "base64encodeddata",
			},
			wantErr: false,
		},
		{
			name: "valid card data",
			data: &AddDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				Type:    "card",
				Payload: "cardnumber:1234567890123456",
			},
			wantErr: false,
		},
		{
			name: "valid otp data",
			data: &AddDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				Type:    "otp",
				Payload: "secretkey",
			},
			wantErr: false,
		},
		{
			name: "empty JWT",
			data: &AddDataFormData{
				JWT:     "",
				Type:    "login",
				Payload: "payload",
			},
			wantErr: true,
		},
		{
			name: "empty type",
			data: &AddDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				Type:    "",
				Payload: "payload",
			},
			wantErr: true,
		},
		{
			name: "empty payload",
			data: &AddDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				Type:    "login",
				Payload: "",
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			data: &AddDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				Type:    "invalid",
				Payload: "payload",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestEditDataFormData_Validate tests the validation of EditDataFormData.
func TestEditDataFormData_Validate(t *testing.T) {
	tests := []struct {
		name    string
		data    *EditDataFormData
		wantErr bool
	}{
		{
			name: "valid data",
			data: &EditDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				ID:      "123e4567-e89b-12d3-a456-426614174000",
				Type:    "login",
				Payload: "updated payload",
			},
			wantErr: false,
		},
		{
			name: "empty JWT",
			data: &EditDataFormData{
				JWT:     "",
				ID:      "123e4567-e89b-12d3-a456-426614174000",
				Type:    "login",
				Payload: "payload",
			},
			wantErr: true,
		},
		{
			name: "empty ID",
			data: &EditDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				ID:      "",
				Type:    "login",
				Payload: "payload",
			},
			wantErr: true,
		},
		{
			name: "empty type",
			data: &EditDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				ID:      "123e4567-e89b-12d3-a456-426614174000",
				Type:    "",
				Payload: "payload",
			},
			wantErr: true,
		},
		{
			name: "empty payload",
			data: &EditDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				ID:      "123e4567-e89b-12d3-a456-426614174000",
				Type:    "login",
				Payload: "",
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			data: &EditDataFormData{
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				ID:      "123e4567-e89b-12d3-a456-426614174000",
				Type:    "invalid",
				Payload: "payload",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestFormDataInterface tests that all form data types implement the FormData interface.
func TestFormDataInterface(t *testing.T) {
	// Test that AuthFormData implements FormData
	var _ FormData = &AuthFormData{}

	// Test that JWTIDFormData implements FormData
	var _ FormData = &JWTIDFormData{}

	// Test that AddDataFormData implements FormData
	var _ FormData = &AddDataFormData{}

	// Test that EditDataFormData implements FormData
	var _ FormData = &EditDataFormData{}

	// Test FormHandler type
	var _ FormHandler[*AuthFormData] = func(*AuthFormData) error { return nil }
	var _ FormHandler[*JWTIDFormData] = func(*JWTIDFormData) error { return nil }
	var _ FormHandler[*AddDataFormData] = func(*AddDataFormData) error { return nil }
	var _ FormHandler[*EditDataFormData] = func(*EditDataFormData) error { return nil }
}
