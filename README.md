# GophKeeper

A secure client-server password manager for storing and syncing private data.

## Features

- **Client-side encryption**: All data is encrypted with AES-256 before being sent to the server
- **TLS security**: All network communications use TLS 1.3
- **JWT authentication**: Secure token-based authentication
- **TUI interface**: Interactive terminal user interface
- **CLI commands**: Command-line interface for automation
- **Cross-platform**: Works on Windows, Linux, and macOS

## Prerequisites

- Go 1.23 or later
- PostgreSQL database
- OpenSSL (for certificate generation)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/AlenaMolokova/gophkeeper.git
cd gophkeeper
```

2. Install dependencies:
```bash
go mod tidy
```

3. Generate TLS certificates:
```bash
mkdir cert
openssl req -x509 -newkey rsa:4096 -keyout cert/server.key -out cert/server.crt -days 365 -nodes
```

## Building

Build the server and client:
```bash
go build -o bin/server cmd/server/main.go
go build -o bin/client cmd/client/main.go
```

## Configuration

### Server Configuration

Set environment variables:
```bash
export DATABASE_DSN="postgres://username:password@localhost:5432/gophkeeper"
export JWT_SECRET="your-secret-key-here"
```

#### Database Connection Pool Configuration

You can configure the database connection pool using environment variables:

```bash
# Pool configuration (optional)
export DB_MAX_CONNS="50"                    # Maximum connections (default: 20)
export DB_MIN_CONNS="10"                    # Minimum connections (default: 5)
export DB_MAX_CONN_LIFETIME="1h"            # Connection lifetime (default: 1h)
export DB_MAX_CONN_IDLE_TIME="30m"          # Idle time (default: 30m)
export DB_HEALTH_CHECK_PERIOD="30s"         # Health check period (default: 1m)
```

For detailed configuration options, see [Pool Configuration Documentation](docs/pool-config.md).

### Database Setup

Create PostgreSQL database and tables:
```sql
CREATE DATABASE gophkeeper;

\c gophkeeper

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    payload BYTEA NOT NULL,
    metadata JSONB,
    timestamp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Usage

### Starting the Server

```bash
./bin/server
```

The server will start on port 50051 with TLS enabled.

### Client Commands

#### Initialize encryption key
```bash
./bin/client init-key
```

#### Register a new user
```bash
./bin/client register --email user@example.com --password secret
```

#### Login
```bash
./bin/client login --email user@example.com --password secret
```

#### Add data
```bash
./bin/client add --token <JWT> --type login --payload "mysecret"
```

#### Get data
```bash
./bin/client get --token <JWT> --id <data_id>
```

#### Edit data
```bash
./bin/client edit --token <JWT> --id <data_id> --type login --payload "newsecret"
```

#### Delete data
```bash
./bin/client delete --token <JWT> --id <data_id>
```

#### Interactive TUI
```bash
./bin/client tui
```

#### Sync operations
```bash
./bin/client sync-upload
./bin/client sync-download
```

#### Show version
```bash
./bin/client version
```

## Testing

### Unit Tests
```bash
go test ./tests/unit/... -v
```

### Integration Tests
```bash
go test ./tests/integration/... -v
```

### Test Coverage
```bash
go test ./... -cover
```

### Race Condition Tests
```bash
go test ./... -race
```

### Static Analysis
```bash
golangci-lint run
```

## Security

- **Client-side encryption**: Data is encrypted with AES-256 before transmission
- **Local key storage**: Encryption keys are stored only on the client machine
- **TLS 1.3**: All network communications are encrypted
- **JWT tokens**: Secure authentication with configurable expiration
- **Input validation**: All inputs are validated and sanitized

## Architecture

```
gophkeeper/
├── api/                 # Protocol buffers and API definitions
├── cmd/                 # Application entry points
│   ├── client/         # Client application
│   └── server/         # Server application
├── internal/           # Internal packages
│   ├── client/         # Client-side logic
│   ├── server/         # Server-side logic
│   └── shared/         # Shared utilities
├── pkg/                # Public packages
└── tests/              # Test files
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License.