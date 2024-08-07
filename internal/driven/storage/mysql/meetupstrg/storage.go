package meetupstrg

import (
	"context"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

type Storage struct {
	sqlClient *sqlx.DB
}

// CancelMeetup implements meetup.MeetupStorage.
func (s *Storage) CancelMeetup(ctx context.Context, meetupID int, cancelledReason string) error {
	panic("unimplemented")
}

// GetMeetup implements meetup.MeetupStorage.
func (s *Storage) GetMeetup(ctx context.Context, meetupID int) (*entity.Meetup, error) {
	panic("unimplemented")
}

// GetMeetups implements meetup.MeetupStorage.
func (s *Storage) GetMeetups(ctx context.Context) ([]entity.Meetup, error) {
	panic("unimplemented")
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

func (s *Storage) SaveMeetup(ctx context.Context, meetup entity.Meetup) error {
	meetupRow := NewMeetupRow(&meetup)
	query := `
		REPLACE INTO meetup (
			id, name, venue_id, event_id, start_ts, end_ts, max_persons, organizer_id, status, created_at, updated_at
		) VALUES (
			:id, :name, :venue_id, :event_id, :start_ts, :end_ts, :max_persons, :organizer_id, :status, :created_at, :updated_at
		)
	`

	_, err := s.sqlClient.NamedExecContext(ctx, query, map[string]interface{}{
		"id":           meetupRow.ID,
		"name":         meetupRow.Name,
		"venue_id":     meetupRow.VenueID,
		"event_id":     meetupRow.EventID,
		"start_ts":     meetupRow.StartTs,
		"end_ts":       meetupRow.EndTs,
		"max_persons":  meetupRow.MaxPersons,
		"organizer_id": meetupRow.OrganizerID,
		"status":       meetupRow.Status,
		"created_at":   meetupRow.CreatedAt,
		"updated_at":   meetupRow.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}
	return nil
}
