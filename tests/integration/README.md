# Integration Tests with Testcontainers

This directory contains integration tests for the GophKeeper application using Testcontainers for complete test isolation.

## Overview

The integration tests use Testcontainers to:
- Spin up isolated PostgreSQL containers for each test
- Generate dynamic TLS certificates
- Eliminate external dependencies and hardcoded paths
- Ensure test reproducibility across different environments

## Architecture

### Testcontainers Infrastructure

- **PostgreSQL Container**: Isolated database instance for each test
- **Dynamic TLS Certificates**: Generated on-the-fly for secure communication
- **Test Environment**: Complete isolation with automatic cleanup

### Key Components

1. **`testutils/containers.go`**: Testcontainers setup and management
2. **`testutils/testutils.go`**: Updated utilities for containerized testing
3. **`auth_test_containers.go`**: Example tests using Testcontainers
4. **`Dockerfile.test`**: Server container for testing

## Running Tests

### Prerequisites

- Docker installed and running
- Go 1.23+
- Testcontainers dependencies installed

### Basic Test Execution

```bash
# Run all integration tests
go test -v ./tests/integration/...

# Run specific test file
go test -v ./tests/integration/auth_test_containers.go

# Run with race detection
go test -race -v ./tests/integration/...
```

### Test Environment Variables

The tests automatically set up the following environment variables:

```bash
DATABASE_DSN=postgres://testuser:testpass@localhost:5432/gophkeeper_test
JWT_SECRET=test-secret-key
CERT_PATH=/tmp/test-certs/server.crt
```

## Benefits of Testcontainers

### 1. Complete Isolation
- Each test runs in its own PostgreSQL container
- No interference between tests
- Clean state for every test run

### 2. No External Dependencies
- No need for local PostgreSQL installation
- No hardcoded certificate paths
- Works on any machine with Docker

### 3. Reproducible Tests
- Same environment across different machines
- Consistent test results
- Easy CI/CD integration

### 4. Dynamic Resource Management
- Automatic container cleanup
- Resource isolation
- No port conflicts

## Migration from Old Tests

### Before (Hardcoded Dependencies)
```go
func TestUserRegistration(t *testing.T) {
    // Hardcoded certificate path
    certPath := "cert/server.crt"
    // External database dependency
    // Manual cleanup required
}
```

### After (Testcontainers)
```go
func TestUserRegistrationWithContainers(t *testing.T) {
    // Isolated environment
    env := testutils.SetupTestEnvironment(t)
    defer env.Cleanup()
    
    // Dynamic certificates
    clients := testutils.SetupTestClients(t, env)
    
    // Test logic...
}
```

## Best Practices

### 1. Always Use Cleanup
```go
env := testutils.SetupTestEnvironment(t)
defer env.Cleanup() // Ensures proper cleanup
```

### 2. Wait for Services
```go
// Give server time to start
time.Sleep(2 * time.Second)
```

### 3. Use Context for Timeouts
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

### 4. Test Error Scenarios
```go
// Test both success and failure cases
// Verify proper error handling
```

## Troubleshooting

### Common Issues

1. **Docker not running**: Ensure Docker is started
2. **Port conflicts**: Tests use isolated containers
3. **Certificate errors**: Certificates are generated dynamically
4. **Database connection issues**: Check container logs

### Debug Mode

Enable verbose logging:
```bash
go test -v -count=1 ./tests/integration/...
```

### Container Logs

Check container logs for debugging:
```go
logs, err := container.Logs(ctx)
if err != nil {
    t.Logf("Container logs: %s", logs)
}
```

## Future Improvements

1. **Parallel Test Execution**: Run tests in parallel with isolated containers
2. **Custom Test Images**: Optimized Docker images for faster startup
3. **Health Checks**: Wait for services to be ready
4. **Test Data Management**: Seeded test data for complex scenarios 