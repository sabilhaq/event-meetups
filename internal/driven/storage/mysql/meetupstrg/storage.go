package meetupstrg

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
func (s *Storage) GetMeetups(ctx context.Context, filter entity.GetMeetupFilter) ([]entity.Meetup, error) {
	var rows MeetupJoinVenueEventRows
	args := []interface{}{}
	queryBuilder := strings.Builder{}

	queryBuilder.WriteString(`
		SELECT
			m.id AS meetup_id,
			m.name AS meetup_name,
			v.id AS venue_id,
			v.name AS venue_name,
			e.id AS event_id,
			e.name AS event_name,
			m.start_ts,
			m.end_ts,
			m.max_persons,
			u.id AS organizer_id,
			u.username AS organizer_username,
			u.email AS organizer_email,
			(SELECT COUNT(*) FROM meetup_user mu WHERE mu.meetup_id = m.id) AS joined_persons_count,
			m.status
		FROM meetup m
		JOIN venue v ON m.venue_id = v.id
		JOIN event e ON m.event_id = e.id
		JOIN user u ON m.organizer_id = u.id
	`)

	conditions := []string{}
	conditions = append(conditions, "m.status = ?")
	args = append(args, "open")

	if filter.EventID != nil {
		conditions = append(conditions, "event_id = ?")
		args = append(args, filter.EventID)
	}

	// Combine conditions with AND
	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY m.start_ts ASC ")

	if filter.Limit != nil {
		queryBuilder.WriteString("LIMIT ?")
		args = append(args, filter.Limit)
	}

	// Finalize the query
	query := queryBuilder.String()

	if err := s.sqlClient.SelectContext(ctx, &rows, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}

	return rows.ToMeetups(), nil
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

func (s *Storage) SaveMeetup(ctx context.Context, meetup entity.Meetup) (int, error) {
	meetupRow := NewMeetupRow(&meetup)
	query := `
		REPLACE INTO meetup (
			id, name, venue_id, event_id, start_ts, end_ts, max_persons, organizer_id, status, created_at, updated_at
		) VALUES (
			:id, :name, :venue_id, :event_id, :start_ts, :end_ts, :max_persons, :organizer_id, :status, :created_at, :updated_at
		)
	`

	res, err := s.sqlClient.NamedExecContext(ctx, query, map[string]interface{}{
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
		return 0, fmt.Errorf("unable to execute query due: %w", err)
	}

	id, err := res.LastInsertId()
	return int(id), err
}

// CountMeetups implements meetup.VenueStorage.
func (s *Storage) CountMeetups(ctx context.Context, venueID, eventID int, startTs, endTs int64) (*int, error) {
	var count int

	query := `
		SELECT COUNT(*) AS meetups_count
		FROM meetup m
		WHERE m.venue_id = ? 
		AND m.event_id = ? 
		AND (
			(m.start_ts <= ? AND m.end_ts >= ?) OR 
			(m.start_ts <= ? AND m.start_ts >= ?) OR 
			(m.start_ts >= ? AND m.start_ts <= ?) 			
		)
	`

	if err := s.sqlClient.GetContext(ctx, &count, query, venueID, eventID, startTs, endTs, startTs, endTs, startTs, endTs); err != nil {
		return nil, fmt.Errorf("unable to find supported event with venue id %d and event id %d: %v", venueID, eventID, err)
	}

	return &count, nil
}
