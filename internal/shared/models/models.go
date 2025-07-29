// Package models provides data structures and type definitions for the GophKeeper application.
// It includes user and data models with conversion utilities for protobuf integration.
//
// The package defines the core domain entities used throughout the application
// and provides type-safe conversions between internal models and protobuf messages.
package models

import (
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// DataType represents the type of stored data.
// It is an alias to the protobuf-generated DataType enum for type safety.
type DataType = clientapi.DataType

// Constants for DataType provide convenient access to data type values.
// These constants ensure type safety when working with different data types
// throughout the application.
const (
	DataTypeUnspecified = clientapi.DataType_DATA_TYPE_UNSPECIFIED
	DataTypeLogin       = clientapi.DataType_DATA_TYPE_LOGIN
	DataTypeText        = clientapi.DataType_DATA_TYPE_TEXT
	DataTypeBinary      = clientapi.DataType_DATA_TYPE_BINARY
	DataTypeCard        = clientapi.DataType_DATA_TYPE_CARD
	DataTypeOTP         = clientapi.DataType_DATA_TYPE_OTP
)

// LoginData represents login credentials metadata.
// It contains structured information about login entries such as usernames,
// URLs, and additional notes for better organization.
type LoginData = clientapi.LoginData

// CardData represents credit card metadata.
// It contains structured information about credit card entries including
// card numbers, holder names, expiry dates, and security codes.
type CardData = clientapi.CardData

// TextData represents text data metadata.
// It contains structured information about text entries such as titles
// and additional notes for better organization.
type TextData = clientapi.TextData

// BinaryData represents binary data metadata.
// It contains structured information about binary file entries including
// filenames, content types, sizes, and additional notes.
type BinaryData = clientapi.BinaryData

// OTPData represents OTP metadata.
// It contains structured information about OTP entries including issuer,
// account, algorithm, digits, period, and additional notes.
type OTPData = clientapi.OTPData

// User represents a user entity in the GophKeeper system.
// It contains the user's unique identifier, email address, and password hash.
// The password is stored as a bcrypt hash for security.
type User struct {
	ID    string // Unique user identifier
	Email string // User's email address (unique)
	Hash  string // Bcrypt hash of the user's password
}

// Data represents a data entry in the GophKeeper system.
// It contains encrypted data with structured metadata and timestamp for synchronization.
// The Payload field contains the encrypted data, while Metadata contains type-specific
// structured information for better organization and validation.
type Data struct {
	ID        string      // Unique data identifier
	Type      DataType    // Type of stored data (login, text, binary, etc.)
	Payload   []byte      // Encrypted data payload
	Metadata  interface{} // Type-specific metadata (LoginData, CardData, etc.)
	Timestamp int64       // Unix timestamp for synchronization
}

// ConvertToProto converts internal Data model to protobuf Data message.
// It handles the conversion of structured metadata to protobuf oneof fields
// based on the data type, ensuring proper serialization for gRPC communication.
//
// The function preserves all data fields and converts the metadata interface
// to the appropriate protobuf message type using type assertions.
func (d *Data) ConvertToProto() *clientapi.Data {
	protoData := &clientapi.Data{
		Id:        d.ID,
		Type:      d.Type,
		Payload:   d.Payload,
		Timestamp: d.Timestamp,
	}

	// Convert metadata based on type
	switch metadata := d.Metadata.(type) {
	case *LoginData:
		protoData.Metadata = &clientapi.Data_LoginData{LoginData: metadata}
	case *CardData:
		protoData.Metadata = &clientapi.Data_CardData{CardData: metadata}
	case *TextData:
		protoData.Metadata = &clientapi.Data_TextData{TextData: metadata}
	case *BinaryData:
		protoData.Metadata = &clientapi.Data_BinaryData{BinaryData: metadata}
	case *OTPData:
		protoData.Metadata = &clientapi.Data_OtpData{OtpData: metadata}
	}

	return protoData
}

// ConvertFromProto converts protobuf Data message to internal Data model.
// It handles the conversion of protobuf oneof metadata fields to internal
// structured types, ensuring proper deserialization from gRPC communication.
//
// The function preserves all data fields and converts the protobuf metadata
// to the appropriate internal type using type assertions.
func ConvertFromProto(protoData *clientapi.Data) *Data {
	data := &Data{
		ID:        protoData.Id,
		Type:      protoData.Type,
		Payload:   protoData.Payload,
		Timestamp: protoData.Timestamp,
	}

	// Convert metadata based on type
	switch metadata := protoData.Metadata.(type) {
	case *clientapi.Data_LoginData:
		data.Metadata = metadata.LoginData
	case *clientapi.Data_CardData:
		data.Metadata = metadata.CardData
	case *clientapi.Data_TextData:
		data.Metadata = metadata.TextData
	case *clientapi.Data_BinaryData:
		data.Metadata = metadata.BinaryData
	case *clientapi.Data_OtpData:
		data.Metadata = metadata.OtpData
	}

	return data
}
