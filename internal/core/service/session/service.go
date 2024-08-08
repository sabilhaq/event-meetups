package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"gopkg.in/validator.v2"
)

var (
	ErrUserNotFound = errors.New("user is not found")
	ErrInvalidCreds = errors.New("invalid username or password")
)

type Service interface {
	// CreateSession is used to login to the system.
	// It returns a JWT token that can be used to access
	// other endpoints named access_token.
	CreateSession(ctx context.Context, username, password string) (*entity.Session, error)
}

type service struct {
	sessionStorage SessionStorage
	userStorage    UserStorage
}

func (s *service) CreateSession(ctx context.Context, username, password string) (*entity.Session, error) {
	// get user
	user, err := s.userStorage.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch user instance due: %w", err)
	}
	if user == nil || user.Password != password {
		return nil, ErrInvalidCreds
	}

	// initiate new session instance
	cfg := entity.SessionConfig{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Exp:      time.Now().Add(time.Hour * 24).Unix(),
	}
	session, err := entity.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize session instance due: %w", err)
	}
	// generate token instance on storage
	accessToken, err := s.sessionStorage.GenerateToken(ctx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to generate token due: %w", err)
	}

	session.AccessToken = accessToken
	return session, nil
}

type ServiceConfig struct {
	SessionStorage SessionStorage `validate:"nonnil"`
	UserStorage    UserStorage    `validate:"nonnil"`
}

func (c ServiceConfig) Validate() error {
	return validator.Validate(c)
}

// NewService returns new instance of service.
func NewService(cfg ServiceConfig) (Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &service{
		sessionStorage: cfg.SessionStorage,
		userStorage:    cfg.UserStorage,
	}
	return s, nil
}
