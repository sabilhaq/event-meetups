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

	// GetVenue returns information about a single venue.
	GetVenue(ctx context.Context, venueID int) (*entity.Venue, error)
}

type service struct {
	venueStorage VenueStorage
	eventStorage EventStorage
}

func (s *service) GetVenues(ctx context.Context, filter entity.GetVenueFilter) ([]entity.Venue, error) {
	venues, err := s.venueStorage.GetVenues(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to get available venues due: %w", err)
	}
	res := make([]entity.Venue, len(venues))
	for i := 0; i < len(venues); i++ {
		supportedEvents, err := s.eventStorage.GetSupportedEvents(ctx, venues[i].ID)
		if err != nil {
			return nil, fmt.Errorf("unable to get supported events due: %w", err)
		}

		res[i] = entity.Venue{
			ID:              venues[i].ID,
			Name:            venues[i].Name,
			OpenDays:        venues[i].OpenDays,
			OpenAt:          venues[i].OpenAt,
			ClosedAt:        venues[i].ClosedAt,
			Timezone:        venues[i].Timezone,
			SupportedEvents: supportedEvents,
		}
	}
	return res, nil
}

func (s *service) GetVenue(ctx context.Context, venueID int) (*entity.Venue, error) {
	// get venue instance from storage
	venue, err := s.venueStorage.GetVenue(ctx, venueID)
	if err != nil {
		return nil, fmt.Errorf("unable to get venue instance due: %w", err)
	}
	if venue == nil {
		return nil, ErrVenueNotFound
	}
	return venue, nil
}

type ServiceConfig struct {
	VenueStorage VenueStorage `validate:"nonnil"`
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
		venueStorage: cfg.VenueStorage,
		eventStorage: cfg.EventStorage,
	}
	return s, nil
}
