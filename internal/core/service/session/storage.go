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
	// GetUserByUsername returns user instance for given username from storage. Returns nil
	// when given username is not found in database.
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}
