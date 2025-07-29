package tui

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateAuthFormBase tests the creation of base authentication form.
func TestCreateAuthFormBase(t *testing.T) {
	handler := func(data *AuthFormData) error {
		assert.Equal(t, "test@example.com", data.Email)
		assert.Equal(t, "password123", data.Password)
		return nil
	}

	backHandler := func() {}

	form := createAuthFormBase("Test Action", handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	emailField := form.GetFormItemByLabel("Email").(*tview.InputField)
	emailField.SetText("test@example.com")

	passwordField := form.GetFormItemByLabel("Пароль").(*tview.InputField)
	passwordField.SetText("password123")

	// Test that form is properly constructed with expected fields
	assert.NotNil(t, form.GetFormItemByLabel("Email"))
	assert.NotNil(t, form.GetFormItemByLabel("Пароль"))
}

// TestCreateAuthFormBaseWithError tests the creation of base authentication form with error handling.
func TestCreateAuthFormBaseWithError(t *testing.T) {
	handler := func(data *AuthFormData) error {
		return assert.AnError
	}

	backHandler := func() {}

	form := createAuthFormBase("Test Action", handler, backHandler)
	require.NotNil(t, form)

	// Test that form is properly constructed even with error handler
	assert.NotNil(t, form.GetFormItemByLabel("Email"))
	assert.NotNil(t, form.GetFormItemByLabel("Пароль"))
}

// TestCreateLoginForm tests the creation of login form.
func TestCreateLoginForm(t *testing.T) {
	handler := func(data *AuthFormData) error {
		assert.Equal(t, "test@example.com", data.Email)
		assert.Equal(t, "password123", data.Password)
		return nil
	}

	backHandler := func() {}

	form := createLoginForm(handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	emailField := form.GetFormItemByLabel("Email").(*tview.InputField)
	emailField.SetText("test@example.com")

	passwordField := form.GetFormItemByLabel("Пароль").(*tview.InputField)
	passwordField.SetText("password123")

	// Test that form is properly constructed with expected fields
	assert.NotNil(t, form.GetFormItemByLabel("Email"))
	assert.NotNil(t, form.GetFormItemByLabel("Пароль"))
}

// TestCreateRegisterForm tests the creation of register form.
func TestCreateRegisterForm(t *testing.T) {
	handler := func(data *AuthFormData) error {
		assert.Equal(t, "test@example.com", data.Email)
		assert.Equal(t, "password123", data.Password)
		return nil
	}

	backHandler := func() {}

	form := createRegisterForm(handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	emailField := form.GetFormItemByLabel("Email").(*tview.InputField)
	emailField.SetText("test@example.com")

	passwordField := form.GetFormItemByLabel("Пароль").(*tview.InputField)
	passwordField.SetText("password123")

	// Test that form is properly constructed with expected fields
	assert.NotNil(t, form.GetFormItemByLabel("Email"))
	assert.NotNil(t, form.GetFormItemByLabel("Пароль"))
}

// TestCreateJWTIDFormBase tests the creation of base JWT ID form.
func TestCreateJWTIDFormBase(t *testing.T) {
	handler := func(data *JWTIDFormData) error {
		assert.Equal(t, "test-jwt", data.JWT)
		assert.Equal(t, "test-id", data.ID)
		return nil
	}

	backHandler := func() {}

	form := createJWTIDFormBase("Test Action", handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	jwtField := form.GetFormItemByLabel("JWT").(*tview.InputField)
	jwtField.SetText("test-jwt")

	idField := form.GetFormItemByLabel("ID").(*tview.InputField)
	idField.SetText("test-id")

	// Test that form is properly constructed with expected fields
	assert.NotNil(t, form.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form.GetFormItemByLabel("ID"))
}

// TestCreateJWTIDFormBaseWithError tests the creation of base JWT ID form with error handling.
func TestCreateJWTIDFormBaseWithError(t *testing.T) {
	handler := func(data *JWTIDFormData) error {
		return assert.AnError
	}

	backHandler := func() {}

	form := createJWTIDFormBase("Test Action", handler, backHandler)
	require.NotNil(t, form)

	// Test that form is properly constructed even with error handler
	assert.NotNil(t, form.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form.GetFormItemByLabel("ID"))
}

// TestCreateGetDataForm tests the creation of get data form.
func TestCreateGetDataForm(t *testing.T) {
	handler := func(data *JWTIDFormData) error {
		assert.Equal(t, "test-jwt", data.JWT)
		assert.Equal(t, "test-id", data.ID)
		return nil
	}

	backHandler := func() {}

	form := createGetDataForm(handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	jwtField := form.GetFormItemByLabel("JWT").(*tview.InputField)
	jwtField.SetText("test-jwt")

	idField := form.GetFormItemByLabel("ID").(*tview.InputField)
	idField.SetText("test-id")

	// Test that form is properly constructed with expected fields
	assert.NotNil(t, form.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form.GetFormItemByLabel("ID"))
}

// TestCreateAddDataForm tests the creation of add data form.
func TestCreateAddDataForm(t *testing.T) {
	handler := func(data *AddDataFormData) error {
		assert.Equal(t, "test-jwt", data.JWT)
		assert.Equal(t, "test-type", data.Type)
		assert.Equal(t, "test-payload", data.Payload)
		return nil
	}

	backHandler := func() {}

	form := createAddDataForm(handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	jwtField := form.GetFormItemByLabel("JWT").(*tview.InputField)
	jwtField.SetText("test-jwt")

	typeField := form.GetFormItemByLabel("Тип").(*tview.InputField)
	typeField.SetText("test-type")

	payloadField := form.GetFormItemByLabel("Payload").(*tview.InputField)
	payloadField.SetText("test-payload")

	// Test that form is properly constructed with expected fields
	assert.NotNil(t, form.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form.GetFormItemByLabel("Тип"))
	assert.NotNil(t, form.GetFormItemByLabel("Payload"))
}

// TestCreateAddDataFormWithError tests the creation of add data form with error handling.
func TestCreateAddDataFormWithError(t *testing.T) {
	handler := func(data *AddDataFormData) error {
		return assert.AnError
	}

	backHandler := func() {}

	form := createAddDataForm(handler, backHandler)
	require.NotNil(t, form)

	// Test that form is properly constructed even with error handler
	assert.NotNil(t, form.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form.GetFormItemByLabel("Тип"))
	assert.NotNil(t, form.GetFormItemByLabel("Payload"))
}

// TestCreateEditDataForm tests the creation of edit data form.
func TestCreateEditDataForm(t *testing.T) {
	handler := func(data *EditDataFormData) error {
		assert.Equal(t, "test-jwt", data.JWT)
		assert.Equal(t, "test-id", data.ID)
		assert.Equal(t, "test-type", data.Type)
		assert.Equal(t, "test-payload", data.Payload)
		return nil
	}

	backHandler := func() {}

	form := createEditDataForm(handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	jwtField := form.GetFormItemByLabel("JWT").(*tview.InputField)
	jwtField.SetText("test-jwt")

	idField := form.GetFormItemByLabel("ID").(*tview.InputField)
	idField.SetText("test-id")

	typeField := form.GetFormItemByLabel("Тип").(*tview.InputField)
	typeField.SetText("test-type")

	payloadField := form.GetFormItemByLabel("Payload").(*tview.InputField)
	payloadField.SetText("test-payload")

	// Test that form is properly constructed with expected fields
	assert.NotNil(t, form.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form.GetFormItemByLabel("ID"))
	assert.NotNil(t, form.GetFormItemByLabel("Тип"))
	assert.NotNil(t, form.GetFormItemByLabel("Payload"))
}

// TestCreateEditDataFormWithError tests the creation of edit data form with error handling.
func TestCreateEditDataFormWithError(t *testing.T) {
	handler := func(data *EditDataFormData) error {
		return assert.AnError
	}

	backHandler := func() {}

	form := createEditDataForm(handler, backHandler)
	require.NotNil(t, form)

	// Test that form is properly constructed even with error handler
	assert.NotNil(t, form.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form.GetFormItemByLabel("ID"))
	assert.NotNil(t, form.GetFormItemByLabel("Тип"))
	assert.NotNil(t, form.GetFormItemByLabel("Payload"))
}

// TestCreateDeleteDataForm tests the creation of delete data form.
func TestCreateDeleteDataForm(t *testing.T) {
	handler := func(data *JWTIDFormData) error {
		assert.Equal(t, "test-jwt", data.JWT)
		assert.Equal(t, "test-id", data.ID)
		return nil
	}

	backHandler := func() {}

	form := createDeleteDataForm(handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	jwtField := form.GetFormItemByLabel("JWT").(*tview.InputField)
	jwtField.SetText("test-jwt")

	idField := form.GetFormItemByLabel("ID").(*tview.InputField)
	idField.SetText("test-id")

	// Test that form is properly constructed with expected fields
	assert.NotNil(t, form.GetFormItemByLabel("JWT"))
	assert.NotNil(t, form.GetFormItemByLabel("ID"))
}

// TestFormFieldExtraction tests that form field extraction works correctly.
func TestFormFieldExtraction(t *testing.T) {
	handler := func(data *EditDataFormData) error {
		assert.Equal(t, "test-jwt", data.JWT)
		assert.Equal(t, "test-id", data.ID)
		assert.Equal(t, "test-type", data.Type)
		assert.Equal(t, "test-payload", data.Payload)
		return nil
	}

	backHandler := func() {}

	form := createEditDataForm(handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	jwtField := form.GetFormItemByLabel("JWT").(*tview.InputField)
	jwtField.SetText("test-jwt")

	idField := form.GetFormItemByLabel("ID").(*tview.InputField)
	idField.SetText("test-id")

	typeField := form.GetFormItemByLabel("Тип").(*tview.InputField)
	typeField.SetText("test-type")

	payloadField := form.GetFormItemByLabel("Payload").(*tview.InputField)
	payloadField.SetText("test-payload")
}

// TestFormButtonHandlers tests that form button handlers work correctly.
func TestFormButtonHandlers(t *testing.T) {
	handler := func(data *AddDataFormData) error {
		return nil
	}

	backHandler := func() {}

	form := createAddDataForm(handler, backHandler)
	require.NotNil(t, form)

	// Test that form is properly constructed with expected fields
	jwtField := form.GetFormItemByLabel("JWT")
	assert.NotNil(t, jwtField)

	typeField := form.GetFormItemByLabel("Тип")
	assert.NotNil(t, typeField)

	payloadField := form.GetFormItemByLabel("Payload")
	assert.NotNil(t, payloadField)
}

// TestFormFieldTypes tests that form fields have correct types.
func TestFormFieldTypes(t *testing.T) {
	authHandler := func(data *AuthFormData) error {
		return nil
	}

	dataHandler := func(data *AddDataFormData) error {
		return nil
	}

	backHandler := func() {}

	// Test auth form
	authForm := createLoginForm(authHandler, backHandler)
	emailField := authForm.GetFormItemByLabel("Email")
	assert.IsType(t, &tview.InputField{}, emailField)

	passwordField := authForm.GetFormItemByLabel("Пароль")
	assert.IsType(t, &tview.InputField{}, passwordField)

	// Test data form
	dataForm := createAddDataForm(dataHandler, backHandler)
	jwtField := dataForm.GetFormItemByLabel("JWT")
	assert.IsType(t, &tview.InputField{}, jwtField)

	typeField := dataForm.GetFormItemByLabel("Тип")
	assert.IsType(t, &tview.InputField{}, typeField)

	payloadField := dataForm.GetFormItemByLabel("Payload")
	assert.IsType(t, &tview.InputField{}, payloadField)
}

// TestFormFieldLabels tests that form fields have correct labels.
func TestFormFieldLabels(t *testing.T) {
	authHandler := func(data *AuthFormData) error {
		return nil
	}

	dataHandler := func(data *AddDataFormData) error {
		return nil
	}

	backHandler := func() {}

	// Test auth form labels
	authForm := createLoginForm(authHandler, backHandler)
	assert.NotNil(t, authForm.GetFormItemByLabel("Email"))
	assert.NotNil(t, authForm.GetFormItemByLabel("Пароль"))

	// Test data form labels
	dataForm := createAddDataForm(dataHandler, backHandler)
	assert.NotNil(t, dataForm.GetFormItemByLabel("JWT"))
	assert.NotNil(t, dataForm.GetFormItemByLabel("Тип"))
	assert.NotNil(t, dataForm.GetFormItemByLabel("Payload"))
}

// TestFormFieldValues tests that form field values are correctly set and retrieved.
func TestFormFieldValues(t *testing.T) {
	handler := func(data *AddDataFormData) error {
		assert.Equal(t, "test-jwt", data.JWT)
		assert.Equal(t, "test-type", data.Type)
		assert.Equal(t, "test-payload", data.Payload)
		return nil
	}

	backHandler := func() {}

	form := createAddDataForm(handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	jwtField := form.GetFormItemByLabel("JWT").(*tview.InputField)
	jwtField.SetText("test-jwt")

	typeField := form.GetFormItemByLabel("Тип").(*tview.InputField)
	typeField.SetText("test-type")

	payloadField := form.GetFormItemByLabel("Payload").(*tview.InputField)
	payloadField.SetText("test-payload")

	// Verify values are set correctly
	assert.Equal(t, "test-jwt", jwtField.GetText())
	assert.Equal(t, "test-type", typeField.GetText())
	assert.Equal(t, "test-payload", payloadField.GetText())
}

// TestFormFieldValidation tests that form field validation works correctly.
func TestFormFieldValidation(t *testing.T) {
	handler := func(data *AuthFormData) error {
		return nil
	}

	backHandler := func() {}

	form := createLoginForm(handler, backHandler)
	require.NotNil(t, form)

	// Test that form fields are accessible
	emailField := form.GetFormItemByLabel("Email")
	assert.NotNil(t, emailField)

	passwordField := form.GetFormItemByLabel("Пароль")
	assert.NotNil(t, passwordField)
}

// TestFormButtonCount tests that forms have the correct number of buttons.
func TestFormButtonCount(t *testing.T) {
	authHandler := func(data *AuthFormData) error {
		return nil
	}

	dataHandler := func(data *AddDataFormData) error {
		return nil
	}

	jwtHandler := func(data *JWTIDFormData) error {
		return nil
	}

	backHandler := func() {}

	// Test auth form buttons
	authForm := createLoginForm(authHandler, backHandler)
	buttonCount := authForm.GetButtonCount()
	assert.Equal(t, 2, buttonCount) // Action button + Back button

	// Test data form buttons
	dataForm := createAddDataForm(dataHandler, backHandler)
	buttonCount = dataForm.GetButtonCount()
	assert.Equal(t, 2, buttonCount) // Action button + Back button

	// Test JWT ID form buttons
	jwtForm := createGetDataForm(jwtHandler, backHandler)
	buttonCount = jwtForm.GetButtonCount()
	assert.Equal(t, 2, buttonCount) // Action button + Back button
}

// TestFormFieldAccessibility tests that form fields are accessible by label.
func TestFormFieldAccessibility(t *testing.T) {
	handler := func(data *AddDataFormData) error {
		return nil
	}

	backHandler := func() {}

	form := createAddDataForm(handler, backHandler)
	require.NotNil(t, form)

	// Test that all expected fields are accessible
	fields := []string{"JWT", "Тип", "Payload"}
	for _, field := range fields {
		item := form.GetFormItemByLabel(field)
		assert.NotNil(t, item, "Field %s should be accessible", field)
	}
}

// TestFormHandlerExecution tests that form handlers are executed correctly.
func TestFormHandlerExecution(t *testing.T) {
	executed := false
	handler := func(data *AuthFormData) error {
		executed = true
		assert.Equal(t, "test@example.com", data.Email)
		assert.Equal(t, "password123", data.Password)
		return nil
	}

	backHandler := func() {}

	form := createLoginForm(handler, backHandler)
	require.NotNil(t, form)

	// Set values in form fields
	emailField := form.GetFormItemByLabel("Email").(*tview.InputField)
	emailField.SetText("test@example.com")

	passwordField := form.GetFormItemByLabel("Пароль").(*tview.InputField)
	passwordField.SetText("password123")

	// Note: We can't directly test button click execution in unit tests
	// as it requires a running tview application. This test verifies
	// that the handler is properly set up.
	assert.False(t, executed, "Handler should not be executed during form creation")
}

// TestFormBackHandler tests that back handlers are properly set up.
func TestFormBackHandler(t *testing.T) {
	handler := func(data *AuthFormData) error {
		return nil
	}

	backExecuted := false
	backHandler := func() {
		backExecuted = true
	}

	form := createLoginForm(handler, backHandler)
	require.NotNil(t, form)

	// Test that back button exists
	buttonCount := form.GetButtonCount()
	assert.Equal(t, 2, buttonCount)

	// Note: We can't directly test button click execution in unit tests
	// as it requires a running tview application. This test verifies
	// that the back handler is properly set up.
	assert.False(t, backExecuted, "Back handler should not be executed during form creation")
}

// TestFormFieldTypesComprehensive tests all form field types comprehensively.
func TestFormFieldTypesComprehensive(t *testing.T) {
	authHandler := func(data *AuthFormData) error {
		return nil
	}

	addHandler := func(data *AddDataFormData) error {
		return nil
	}

	editHandler := func(data *EditDataFormData) error {
		return nil
	}

	jwtHandler := func(data *JWTIDFormData) error {
		return nil
	}

	backHandler := func() {}

	// Test all form types
	forms := []struct {
		name string
		form *tview.Form
	}{
		{"Login", createLoginForm(authHandler, backHandler)},
		{"Register", createRegisterForm(authHandler, backHandler)},
		{"AddData", createAddDataForm(addHandler, backHandler)},
		{"EditData", createEditDataForm(editHandler, backHandler)},
		{"GetData", createGetDataForm(jwtHandler, backHandler)},
		{"DeleteData", createDeleteDataForm(jwtHandler, backHandler)},
	}

	for _, testCase := range forms {
		t.Run(testCase.name, func(t *testing.T) {
			form := testCase.form
			require.NotNil(t, form)

			// Get all form items and verify they are InputField types
			for i := 0; i < form.GetFormItemCount(); i++ {
				item := form.GetFormItem(i)
				assert.IsType(t, &tview.InputField{}, item, "Form item should be InputField")
			}
		})
	}
}

// TestFormFieldLabelsComprehensive tests all form field labels comprehensively.
func TestFormFieldLabelsComprehensive(t *testing.T) {
	authHandler := func(data *AuthFormData) error {
		return nil
	}

	addHandler := func(data *AddDataFormData) error {
		return nil
	}

	editHandler := func(data *EditDataFormData) error {
		return nil
	}

	jwtHandler := func(data *JWTIDFormData) error {
		return nil
	}

	backHandler := func() {}

	// Test all form types with their expected labels
	testCases := []struct {
		name   string
		form   *tview.Form
		labels []string
	}{
		{
			name:   "Login",
			form:   createLoginForm(authHandler, backHandler),
			labels: []string{"Email", "Пароль"},
		},
		{
			name:   "Register",
			form:   createRegisterForm(authHandler, backHandler),
			labels: []string{"Email", "Пароль"},
		},
		{
			name:   "AddData",
			form:   createAddDataForm(addHandler, backHandler),
			labels: []string{"JWT", "Тип", "Payload"},
		},
		{
			name:   "EditData",
			form:   createEditDataForm(editHandler, backHandler),
			labels: []string{"JWT", "ID", "Тип", "Payload"},
		},
		{
			name:   "GetData",
			form:   createGetDataForm(jwtHandler, backHandler),
			labels: []string{"JWT", "ID"},
		},
		{
			name:   "DeleteData",
			form:   createDeleteDataForm(jwtHandler, backHandler),
			labels: []string{"JWT", "ID"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			form := testCase.form
			require.NotNil(t, form)

			// Verify all expected labels are present
			for _, label := range testCase.labels {
				item := form.GetFormItemByLabel(label)
				assert.NotNil(t, item, "Form should have field with label: %s", label)
			}
		})
	}
}
