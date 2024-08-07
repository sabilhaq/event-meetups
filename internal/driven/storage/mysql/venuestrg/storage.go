package venuestrg

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/shared"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

type Storage struct {
	sqlClient *sqlx.DB
}

type Config struct {
	SQLClient *sqlx.DB `validate:"nonnil"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &Storage{sqlClient: cfg.SQLClient}
	return s, nil
}

func (s *Storage) GetVenues(ctx context.Context, filter entity.GetVenueFilter) ([]entity.Venue, error) {
	var rows shared.VenueEventRows
	args := []interface{}{}

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT
			v.id AS venue_id,
			v.name AS venue_name,
			v.open_days AS venue_open_days,
			v.open_at AS venue_open_at,
			v.closed_at AS venue_closed_at,
			v.timezone AS venue_timezone,
			e.id AS event_id,
			e.name AS event_name,
			ve.meetups_capacity AS meetups_capacity
		FROM venue v
		JOIN venue_event ve ON v.id = ve.venue_id
		JOIN event e ON ve.event_id = e.id
	`)

	conditions := []string{}

	if filter.EventID != nil {
		conditions = append(conditions, "event_id = ?")
		args = append(args, filter.EventID)
	}

	if filter.MeetupStartTS != nil {
		conditions = append(conditions, "open_at < ?")
		args = append(args, filter.MeetupStartTS)
	}

	if filter.MeetupEndTS != nil {
		conditions = append(conditions, "closed_at < ?")
		args = append(args, filter.MeetupEndTS)
	}

	// Combine conditions with AND
	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY v.id ASC")

	// Finalize the query
	query := queryBuilder.String()

	if err := s.sqlClient.SelectContext(ctx, &rows, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}

	return rows.ToVenues(), nil
}

func (s *Storage) GetVenue(ctx context.Context, venueID int) (*entity.Venue, error) {
	var venues shared.VenueEventRows
	query := `
		SELECT
			v.id AS venue_id,
			v.name AS venue_name,
			v.open_days AS venue_open_days,
			v.open_at AS venue_open_at,
			v.closed_at AS venue_closed_at,
			v.timezone AS venue_timezone,
			e.id AS event_id,
			e.name AS event_name,
			ve.meetups_capacity AS meetups_capacity
		FROM venue v
		JOIN venue_event ve ON v.id = ve.venue_id
		JOIN event e ON ve.event_id = e.id
		WHERE v.id = ?
	`

	if err := s.sqlClient.SelectContext(ctx, &venues, query, venueID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to find venue with id %d: %v", venueID, err)
	}

	if len(venues) == 0 {
		return nil, nil
	}

	return &venues.ToVenues()[0], nil
}

// IsEventSupported implements meetup.VenueStorage.
func (s *Storage) IsEventSupported(ctx context.Context, venueID int, eventID int) (bool, error) {
	panic("unimplemented")
}
