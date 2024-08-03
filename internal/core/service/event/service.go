package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"gopkg.in/validator.v2"
)

var (
	ErrEventNotFound = errors.New("event is not found")
)

type Service interface {
	GetEvents(ctx context.Context) ([]entity.Event, error)
}

type service struct {
	eventStorage EventStorage
}

func (s *service) GetEvents(ctx context.Context) ([]entity.Event, error) {
	events, err := s.eventStorage.GetEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get available events due: %w", err)
	}
	return events, nil
}

type ServiceConfig struct {
	EventStorage EventStorage `validate:"nonnil"`
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
		eventStorage: cfg.EventStorage,
	}
	return s, nil
}
