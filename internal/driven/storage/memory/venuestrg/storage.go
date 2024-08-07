package venuestrg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"gopkg.in/validator.v2"
)

type Storage struct {
	data map[int]entity.Venue
}

// GetVenue implements venue.VenueStorage.
func (s *Storage) GetVenue(ctx context.Context, venueID int) (*entity.Venue, error) {
	v, exist := s.data[venueID]
	if !exist {
		// if item is not found, returns nil as expected by venue interface
		return nil, nil
	}

	return &v, nil
}

// GetVenues implements venue.VenueStorage.
func (s *Storage) GetVenues(ctx context.Context, filter entity.GetVenueFilter) ([]entity.Venue, error) {
	var venues []entity.Venue
	for _, venue := range s.data {
		venues = append(venues, venue)
	}
	return venues, nil
}

type Config struct {
	VenueData []byte `validate:"nonzero"`
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
	// parse venue data
	var rows []venueRow
	err = json.Unmarshal(cfg.VenueData, &rows)
	if err != nil {
		return nil, fmt.Errorf("unable to parse venue data due: %w", err)
	}
	data := map[int]entity.Venue{}
	for _, venueRow := range rows {
		venue := venueRow.toVenue()
		data[venue.ID] = venue
	}
	return &Storage{data: data}, nil
}
