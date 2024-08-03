package userstrg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"gopkg.in/validator.v2"
)

type Storage struct {
	data map[int]entity.User
}

// GetUser implements session.UserStorage.
func (s *Storage) GetUser(ctx context.Context, username string, password string) (*entity.User, error) {
	for _, user := range s.data {
		if user.Username == username && user.Password == password {
			return &user, nil
		}
	}

	return nil, nil
}

// GetUsers implements user.UserStorage.
func (s *Storage) GetUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	for _, user := range s.data {
		users = append(users, user)
	}
	return users, nil
}

type Config struct {
	UserData []byte `validate:"nonzero"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	// validate config
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	// parse user data
	var rows []userRow
	err = json.Unmarshal(cfg.UserData, &rows)
	if err != nil {
		return nil, fmt.Errorf("unable to parse user data due: %w", err)
	}
	data := map[int]entity.User{}
	for _, userRow := range rows {
		user := userRow.toUser()
		data[user.ID] = user
	}
	return &Storage{data: data}, nil
}
