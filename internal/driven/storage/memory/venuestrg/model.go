package venuestrg

import "github.com/Haraj-backend/hex-monscape/internal/core/entity"

type venueRow struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (r venueRow) toVenue() entity.Venue {
	return entity.Venue{
		ID:   r.ID,
		Name: r.Name,
	}
}
