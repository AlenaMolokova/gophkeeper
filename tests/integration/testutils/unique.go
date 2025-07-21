package testutils

import (
	"math/rand"
	"strconv"
)

func UniqueEmail(base string) string {
	return base + strconv.Itoa(rand.Intn(1_000_000)) + "@example.com"
}
