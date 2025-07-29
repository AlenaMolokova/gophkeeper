// Package applogic provides high-level business logic functions for the GophKeeper client.
// It encapsulates common operations like user authentication, data encryption/decryption,
// and gRPC communication with the server.
//
// The package serves as an abstraction layer between the UI and the underlying
// gRPC clients, providing a clean API for client applications.
package applogic

import (
	"fmt"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// ConvertToDataType converts a string representation of data type to the corresponding DataType enum.
// It provides a centralized way to handle data type conversion with proper error handling.
//
// The function follows the Single Responsibility Principle by focusing solely on
// type conversion logic and can be easily tested and reused across the application.
func ConvertToDataType(dataType string) (clientapi.DataType, error) {
	switch dataType {
	case "login":
		return clientapi.DataType_DATA_TYPE_LOGIN, nil
	case "text":
		return clientapi.DataType_DATA_TYPE_TEXT, nil
	case "binary":
		return clientapi.DataType_DATA_TYPE_BINARY, nil
	case "card":
		return clientapi.DataType_DATA_TYPE_CARD, nil
	case "otp":
		return clientapi.DataType_DATA_TYPE_OTP, nil
	default:
		return clientapi.DataType_DATA_TYPE_UNSPECIFIED, fmt.Errorf("unknown data type: %s", dataType)
	}
}
