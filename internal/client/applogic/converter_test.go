package applogic

import (
	"testing"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/stretchr/testify/assert"
)

func TestConvertToDataType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected clientapi.DataType
		hasError bool
	}{
		{
			name:     "valid login type",
			input:    "login",
			expected: clientapi.DataType_DATA_TYPE_LOGIN,
			hasError: false,
		},
		{
			name:     "valid text type",
			input:    "text",
			expected: clientapi.DataType_DATA_TYPE_TEXT,
			hasError: false,
		},
		{
			name:     "valid binary type",
			input:    "binary",
			expected: clientapi.DataType_DATA_TYPE_BINARY,
			hasError: false,
		},
		{
			name:     "valid card type",
			input:    "card",
			expected: clientapi.DataType_DATA_TYPE_CARD,
			hasError: false,
		},
		{
			name:     "valid otp type",
			input:    "otp",
			expected: clientapi.DataType_DATA_TYPE_OTP,
			hasError: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: clientapi.DataType_DATA_TYPE_UNSPECIFIED,
			hasError: true,
		},
		{
			name:     "unknown type",
			input:    "unknown",
			expected: clientapi.DataType_DATA_TYPE_UNSPECIFIED,
			hasError: true,
		},
		{
			name:     "case sensitive login",
			input:    "Login",
			expected: clientapi.DataType_DATA_TYPE_UNSPECIFIED,
			hasError: true,
		},
		{
			name:     "whitespace padded",
			input:    " login ",
			expected: clientapi.DataType_DATA_TYPE_UNSPECIFIED,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToDataType(tt.input)

			if tt.hasError {
				assert.Error(t, err)
				assert.Equal(t, tt.expected, result)
				assert.Contains(t, err.Error(), "unknown data type")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
