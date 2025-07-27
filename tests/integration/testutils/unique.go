// Package testutils provides utility functions for integration testing of the GophKeeper application.
// This file contains helper functions for generating unique test data to avoid conflicts
// during parallel test execution.
package testutils

import (
	"math/rand"
	"strconv"
)

// UniqueEmail generates a unique email address for testing purposes.
// It appends a random number to the base string to ensure uniqueness
// across multiple test runs and parallel test execution.
func UniqueEmail(base string) string {
	return base + strconv.Itoa(rand.Intn(1_000_000)) + "@example.com"
}
