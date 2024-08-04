package venuestrg

import (
	"strconv"
	"strings"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type VenueRow struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	OpenDays string `db:"open_days"`
	OpenAt   string `db:"open_at"`
	ClosedAt string `db:"closed_at"`
	Timezone string `db:"timezone"`
}

type VenueRows []VenueRow

func (r *VenueRow) ToVenue() *entity.Venue {
	// Parse openDays from comma-separated string to []int
	var ods []int
	for _, day := range strings.Split(r.OpenDays, ",") {
		dayInt, _ := strconv.Atoi(day)
		ods = append(ods, dayInt)
	}

	return &entity.Venue{
		ID:       r.ID,
		Name:     r.Name,
		OpenDays: ods,
		OpenAt:   r.OpenAt,
		ClosedAt: r.ClosedAt,
		Timezone: r.Timezone,
	}
}

func (r VenueRows) ToVenues() []entity.Venue {
	var venues []entity.Venue
	for _, row := range r {
		venues = append(venues, *row.ToVenue())
	}
	return venues
}
