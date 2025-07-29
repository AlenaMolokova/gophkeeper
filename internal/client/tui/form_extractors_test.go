package tui

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExtractInputFieldText tests the extractInputFieldText function.
func TestExtractInputFieldText(t *testing.T) {
	form := tview.NewForm()
	form.AddInputField("TestField", "test-value", 20, nil, nil)

	// Test successful extraction
	value, err := extractInputFieldText(form, "TestField")
	require.NoError(t, err)
	assert.Equal(t, "test-value", value)

	// Test field not found
	_, err = extractInputFieldText(form, "NonExistentField")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "form field 'NonExistentField' not found")
}

// TestExtractPasswordFieldText tests the extractPasswordFieldText function.
func TestExtractPasswordFieldText(t *testing.T) {
	form := tview.NewForm()
	form.AddPasswordField("PasswordField", "secret-password", 20, '*', nil)

	// Test successful extraction
	value, err := extractPasswordFieldText(form, "PasswordField")
	require.NoError(t, err)
	assert.Equal(t, "secret-password", value)

	// Test field not found
	_, err = extractPasswordFieldText(form, "NonExistentField")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "form field 'NonExistentField' not found")
}

// TestExtractAuthFormData tests the ExtractAuthFormData function.
func TestExtractAuthFormData(t *testing.T) {
	form := tview.NewForm()
	form.AddInputField("Email", "test@example.com", 30, nil, nil)
	form.AddPasswordField("Пароль", "password123", 30, '*', nil)

	// Test successful extraction
	data, err := ExtractAuthFormData(form)
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", data.Email)
	assert.Equal(t, "password123", data.Password)

	// Test validation failure - empty email
	form = tview.NewForm()
	form.AddInputField("Email", "", 30, nil, nil)
	form.AddPasswordField("Пароль", "password123", 30, '*', nil)

	_, err = ExtractAuthFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test validation failure - empty password
	form = tview.NewForm()
	form.AddInputField("Email", "test@example.com", 30, nil, nil)
	form.AddPasswordField("Пароль", "", 30, '*', nil)

	_, err = ExtractAuthFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

// TestExtractJWTIDFormData tests the ExtractJWTIDFormData function.
func TestExtractJWTIDFormData(t *testing.T) {
	form := tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("ID", "data-id", 36, nil, nil)

	// Test successful extraction
	data, err := ExtractJWTIDFormData(form)
	require.NoError(t, err)
	assert.Equal(t, "jwt-token", data.JWT)
	assert.Equal(t, "data-id", data.ID)

	// Test validation failure - empty JWT
	form = tview.NewForm()
	form.AddInputField("JWT", "", 60, nil, nil)
	form.AddInputField("ID", "data-id", 36, nil, nil)

	_, err = ExtractJWTIDFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test validation failure - empty ID
	form = tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("ID", "", 36, nil, nil)

	_, err = ExtractJWTIDFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

// TestExtractAddDataFormData tests the ExtractAddDataFormData function.
func TestExtractAddDataFormData(t *testing.T) {
	form := tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("Тип", "login", 10, nil, nil)
	form.AddInputField("Payload", "test-payload", 60, nil, nil)

	// Test successful extraction
	data, err := ExtractAddDataFormData(form)
	require.NoError(t, err)
	assert.Equal(t, "jwt-token", data.JWT)
	assert.Equal(t, "login", data.Type)
	assert.Equal(t, "test-payload", data.Payload)

	// Test validation failure - empty JWT
	form = tview.NewForm()
	form.AddInputField("JWT", "", 60, nil, nil)
	form.AddInputField("Тип", "login", 10, nil, nil)
	form.AddInputField("Payload", "test-payload", 60, nil, nil)

	_, err = ExtractAddDataFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test validation failure - empty type
	form = tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("Тип", "", 10, nil, nil)
	form.AddInputField("Payload", "test-payload", 60, nil, nil)

	_, err = ExtractAddDataFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test validation failure - empty payload
	form = tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("Тип", "login", 10, nil, nil)
	form.AddInputField("Payload", "", 60, nil, nil)

	_, err = ExtractAddDataFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test validation failure - invalid type
	form = tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("Тип", "invalid-type", 10, nil, nil)
	form.AddInputField("Payload", "test-payload", 60, nil, nil)

	_, err = ExtractAddDataFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

// TestExtractEditDataFormData tests the ExtractEditDataFormData function.
func TestExtractEditDataFormData(t *testing.T) {
	form := tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("ID", "data-id", 36, nil, nil)
	form.AddInputField("Тип", "login", 10, nil, nil)
	form.AddInputField("Payload", "test-payload", 60, nil, nil)

	// Test successful extraction
	data, err := ExtractEditDataFormData(form)
	require.NoError(t, err)
	assert.Equal(t, "jwt-token", data.JWT)
	assert.Equal(t, "data-id", data.ID)
	assert.Equal(t, "login", data.Type)
	assert.Equal(t, "test-payload", data.Payload)

	// Test validation failure - empty JWT
	form = tview.NewForm()
	form.AddInputField("JWT", "", 60, nil, nil)
	form.AddInputField("ID", "data-id", 36, nil, nil)
	form.AddInputField("Тип", "login", 10, nil, nil)
	form.AddInputField("Payload", "test-payload", 60, nil, nil)

	_, err = ExtractEditDataFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test validation failure - empty ID
	form = tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("ID", "", 36, nil, nil)
	form.AddInputField("Тип", "login", 10, nil, nil)
	form.AddInputField("Payload", "test-payload", 60, nil, nil)

	_, err = ExtractEditDataFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test validation failure - empty type
	form = tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("ID", "data-id", 36, nil, nil)
	form.AddInputField("Тип", "", 10, nil, nil)
	form.AddInputField("Payload", "test-payload", 60, nil, nil)

	_, err = ExtractEditDataFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test validation failure - empty payload
	form = tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("ID", "data-id", 36, nil, nil)
	form.AddInputField("Тип", "login", 10, nil, nil)
	form.AddInputField("Payload", "", 60, nil, nil)

	_, err = ExtractEditDataFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test validation failure - invalid type
	form = tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)
	form.AddInputField("ID", "data-id", 36, nil, nil)
	form.AddInputField("Тип", "invalid-type", 10, nil, nil)
	form.AddInputField("Payload", "test-payload", 60, nil, nil)

	_, err = ExtractEditDataFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

// TestExtractFormDataWithMissingFields tests extraction with missing form fields.
func TestExtractFormDataWithMissingFields(t *testing.T) {
	// Test missing email field
	form := tview.NewForm()
	form.AddPasswordField("Пароль", "password123", 30, '*', nil)

	_, err := ExtractAuthFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to extract email")

	// Test missing password field
	form = tview.NewForm()
	form.AddInputField("Email", "test@example.com", 30, nil, nil)

	_, err = ExtractAuthFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to extract password")

	// Test missing JWT field
	form = tview.NewForm()
	form.AddInputField("ID", "data-id", 36, nil, nil)

	_, err = ExtractJWTIDFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to extract JWT")

	// Test missing ID field
	form = tview.NewForm()
	form.AddInputField("JWT", "jwt-token", 60, nil, nil)

	_, err = ExtractJWTIDFormData(form)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to extract ID")
}

// TestExtractFormDataWithAllValidTypes tests extraction with all valid data types.
func TestExtractFormDataWithAllValidTypes(t *testing.T) {
	validTypes := []string{"login", "text", "binary", "card", "otp"}

	for _, dataType := range validTypes {
		t.Run(dataType, func(t *testing.T) {
			form := tview.NewForm()
			form.AddInputField("JWT", "jwt-token", 60, nil, nil)
			form.AddInputField("Тип", dataType, 10, nil, nil)
			form.AddInputField("Payload", "test-payload", 60, nil, nil)

			data, err := ExtractAddDataFormData(form)
			require.NoError(t, err)
			assert.Equal(t, dataType, data.Type)
		})
	}
}
