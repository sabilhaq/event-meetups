package eventstrg

import (
	"context"
	"database/sql"
	"fmt"

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

func (s *Storage) GetEvents(ctx context.Context) ([]entity.Event, error) {
	var rows shared.EventRows
	args := []interface{}{}
	query := `
		SELECT
			id,
			name
		FROM event
	`

	if err := s.sqlClient.SelectContext(ctx, &rows, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}

	return rows.ToEvents(), nil
}

func (s *Storage) GetEvent(ctx context.Context, eventID int) (*entity.Event, error) {
	var event EventRow
	query := `
		SELECT
			id,
			name
		FROM event
		WHERE eventID = ?
	`

	if err := s.sqlClient.GetContext(ctx, &event, query, eventID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to find event with id %d: %v", eventID, err)
	}

	return event.ToEvent(), nil
}

func (s *Storage) GetSupportedEvents(ctx context.Context, venueID int) ([]entity.SupportedEvent, error) {
	var rows shared.SupportedEventRows
	query := `
		SELECT
			e.id AS id,
			e.name AS name,
			ve.meetups_capacity AS meetups_capacity
		FROM venue_event ve
		LEFT JOIN event e ON ve.event_id = e.id
		WHERE ve.venue_id = ?
	`

	if err := s.sqlClient.SelectContext(ctx, &rows, query, venueID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}

	return rows.ToSupportedEvents(), nil
}
