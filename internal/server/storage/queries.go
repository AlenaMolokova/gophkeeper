// Package storage provides database storage functionality for the GophKeeper server.
// This file contains SQL queries for user and data operations, following the rule
// from pravila.md to avoid inlining SQL queries in the code.
package storage

// SQL queries for user operations.
const (
	// SaveUserQuery inserts a new user and returns the user ID.
	SaveUserQuery = `
		INSERT INTO users (email, hash) 
		VALUES ($1, $2) 
		RETURNING id`

	// FindUserByEmailQuery retrieves a user by email.
	FindUserByEmailQuery = `
		SELECT id, email, hash 
		FROM users 
		WHERE email = $1`
)

// SQL queries for data operations.
const (
	// SaveDataQuery inserts new data and returns the data ID.
	SaveDataQuery = `
		INSERT INTO data (user_id, type, payload, metadata, timestamp) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`

	// FindDataByIDQuery retrieves data by ID for a specific user.
	FindDataByIDQuery = `
		SELECT id, user_id, type, payload, metadata, timestamp 
		FROM data 
		WHERE id = $1 AND user_id = $2`

	// EditDataQuery updates existing data.
	EditDataQuery = `
		UPDATE data 
		SET type = $3, payload = $4, metadata = $5, timestamp = $6 
		WHERE id = $1 AND user_id = $2 
		RETURNING id`

	// DeleteDataQuery removes data by ID for a specific user.
	DeleteDataQuery = `
		DELETE FROM data 
		WHERE id = $1 AND user_id = $2`
)
