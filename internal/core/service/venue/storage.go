package venue

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type VenueStorage interface {
	// GetVenues returns list of venue available in the system.
	// Returns nil when there is no venues available.
	GetVenues(ctx context.Context, filter entity.GetVenueFilter) ([]entity.Venue, error)
}

type EventStorage interface {
	// GetEvents returns list of event that supported by the venue.
	// Returns nil when there is no events available.
	GetSupportedEvents(ctx context.Context, venueID int) ([]entity.SupportedEvent, error)
}
