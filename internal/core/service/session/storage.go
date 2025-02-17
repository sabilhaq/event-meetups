package session

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type SessionStorage interface {
	// GenerateToken is used for generate jwt in storage.
	GenerateToken(ctx context.Context, userID int) (string, error)
}

type UserStorage interface {
	// GetUser returns game instance for given username and password from storage. Returns nil
	// when given username and password is not found in database.
	GetUser(ctx context.Context, username, password string) (*entity.User, error)
}
