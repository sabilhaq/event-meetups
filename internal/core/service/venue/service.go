package venue

import (
	"context"
	"errors"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"gopkg.in/validator.v2"
)

var (
	ErrVenueNotFound = errors.New("venue is not found")
)

type Service interface {
	// GetVenues returns all venues available in the system.
	GetVenues(ctx context.Context, filter entity.GetVenueFilter) ([]entity.Venue, error)
}

type service struct {
	venueStorage VenueStorage
}

func (s *service) GetVenues(ctx context.Context, filter entity.GetVenueFilter) ([]entity.Venue, error) {
	venues, err := s.venueStorage.GetVenues(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to get available venues due: %w", err)
	}
	return venues, nil
}

type ServiceConfig struct {
	VenueStorage VenueStorage `validate:"nonnil"`
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
		venueStorage: cfg.VenueStorage,
	}
	return s, nil
}
