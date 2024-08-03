package token

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/validator.v2"
)

var (
	secretKey = []byte("event-maker")
)

type Storage struct{}

// TODO: change description
// GenerateToken implements token.SessionStorage.
func (s *Storage) GenerateToken(ctx context.Context, userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

type Config struct{}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	// validate config
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &Storage{}, nil
}
