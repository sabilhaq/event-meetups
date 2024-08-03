package venue

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type VenueStorage interface {
	// GetVenues returns list of venue available in the system.
	// Returns nil when there is no venues available.
	GetVenues(ctx context.Context) ([]entity.Venue, error)
}
