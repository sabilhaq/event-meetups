package event

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type EventStorage interface {
	// GetEvents returns list of event that supported by the system.
	// Returns nil when there is no events available.
	GetEvents(ctx context.Context) ([]entity.Event, error)
}
