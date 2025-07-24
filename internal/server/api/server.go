package api

import (
	"github.com/AlenaMolokova/gophkeeper/internal/server/auth"
	"github.com/AlenaMolokova/gophkeeper/internal/server/storage"
	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
)

// Server implements the GophKeeper gRPC service.
type Server struct {
	clientapi.UnimplementedGophKeeperServer
	auth      *auth.Auth
	storage   storage.Storage
	jwtSecret []byte
}

// NewServer creates a new Server instance.
func NewServer(auth *auth.Auth, storage storage.Storage, jwtSecret []byte) *Server {
	return &Server{
		auth:      auth,
		storage:   storage,
		jwtSecret: jwtSecret,
	}
}
