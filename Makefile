# Makefile for GophKeeper project
# Provides convenient commands for building, testing, and development

.PHONY: help build test test-unit test-integration test-parallel test-all clean docker-build docker-run

# Default target
help:
	@echo "Available commands:"
	@echo "  build          - Build server and client binaries"
	@echo "  test-unit      - Run unit tests"
	@echo "  test-integration - Run integration tests with Testcontainers"
	@echo "  test-parallel  - Run parallel integration tests"
	@echo "  test-all       - Run all tests (unit + integration + parallel)"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race detection"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image for testing"
	@echo "  docker-run     - Run application in Docker"

# Build targets
build:
	@echo "Building server and client..."
	go build -o bin/server cmd/server/main.go
	go build -o bin/client cmd/client/main.go
	@echo "Build complete!"

# Test targets
test-unit:
	@echo "Running unit tests..."
	go test -v -race ./tests/unit/...

test-integration:
	@echo "Running integration tests with Testcontainers..."
	@echo "Make sure Docker is running!"
	go test -v -race -timeout 10m ./tests/integration/...

test-parallel:
	@echo "Running parallel integration tests..."
	@echo "Make sure Docker is running!"
	go test -v -race -timeout 15m ./tests/integration/parallel_test.go

test-all: test-unit test-integration test-parallel
	@echo "All tests completed!"

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out -covermode=atomic ./tests/integration/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	@echo "Running tests with race detection..."
	go test -v -race -timeout 10m ./tests/integration/...

# Development targets
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Docker targets
docker-build:
	@echo "Building Docker image for testing..."
	docker build -f tests/integration/Dockerfile.test -t gophkeeper-test .

docker-run:
	@echo "Running application in Docker..."
	docker run --rm -p 50051:50051 gophkeeper-test

# CI/CD targets
ci-test:
	@echo "Running CI test suite..."
	go test -v -race -timeout 10m ./tests/unit/...
	go test -v -race -timeout 15m ./tests/integration/...
	go test -v -race -timeout 20m ./tests/integration/parallel_test.go

# Performance testing
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./tests/integration/parallel_test.go

# Security testing
security-test:
	@echo "Running security tests..."
	go test -v -timeout 10m ./tests/integration/security_test.go

# Database setup for local development
setup-db:
	@echo "Setting up local PostgreSQL database..."
	@echo "This requires PostgreSQL to be installed locally"
	@echo "Consider using Docker: docker run --name postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15"

# Generate mocks
generate-mocks:
	@echo "Generating mocks..."
	go generate ./...

# Linting
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
format:
	@echo "Formatting code..."
	go fmt ./...
	go vet ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy 