package models

// User represents a user entity in the GophKeeper system.
// It contains the user's unique identifier, email address, and password hash.
type User struct {
	ID    string
	Email string
	Hash  string
}

// Data represents a data entry in the GophKeeper system.
// It contains encrypted data with metadata and timestamp for synchronization.
type Data struct {
	ID        string
	Type      string
	Payload   []byte
	Metadata  map[string]string
	Timestamp int64
}
