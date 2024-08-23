package meetupstrg

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
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

// GetMeetup implements meetup.MeetupStorage.
func (s *Storage) GetMeetup(ctx context.Context, meetupID, userID int) (*entity.Meetup, bool, error) {
	var meetup MeetupJoinVenueEventUserRow
	query := `
		SELECT
			m.id,
			m.name,
			v.id AS "venue.id",
			v.name AS "venue.name",
			e.id AS "event.id",
			e.name AS "event.name",
			m.start_ts,
			m.end_ts,
			m.max_persons,
			u.id as "organizer.id", 
			u.username as "organizer.username", 
			u.email as "organizer.email",
			(SELECT COUNT(*) FROM meetup_user mu WHERE mu.meetup_id = m.id) AS joined_persons_count,
			EXISTS (SELECT 1 FROM meetup_user mu WHERE mu.meetup_id = m.id AND mu.user_id = ?) AS is_joined,
			(u.id = ? OR EXISTS (SELECT 1 FROM meetup_user mu WHERE mu.meetup_id = m.id AND mu.user_id = ?)) AS is_organizer_or_participant,
			m.status,
			m.cancelled_reason,
			m.cancelled_at
		FROM meetup m
		JOIN venue v ON m.venue_id = v.id
		JOIN event e ON m.event_id = e.id
		JOIN user u ON m.organizer_id = u.id
		WHERE m.id = ?
	`

	if err := s.sqlClient.GetContext(ctx, &meetup, query, userID, userID, userID, meetupID); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("unable to find meetup with id %d: %v", meetupID, err)
	}

	if meetup.IsOrganizerOrParticipant {
		var joinedPersons []JoinedPerson

		// Define the query for joined persons
		joinedPersonsQuery := `
			SELECT 
				u.id, u.username, u.email, mu.joined_at 
			FROM meetup_user mu 
			JOIN user u ON mu.user_id = u.id 
			WHERE mu.meetup_id = ?
		`

		// Execute the query for joined persons
		if err := s.sqlClient.SelectContext(ctx, &joinedPersons, joinedPersonsQuery, meetupID); err != nil {
			if err == sql.ErrNoRows {
				return nil, false, nil
			}
			return nil, false, fmt.Errorf("unable to execute query due: %w", err)
		}

		// Ensure that `joined_persons` appears as an empty array if there are no joined persons
		if len(joinedPersons) == 0 {
			joinedPersons = []JoinedPerson{}
		}
		meetup.JoinedPersons = joinedPersons
	}

	return meetup.ToMeetup(), meetup.IsOrganizerOrParticipant, nil
}

// CancelMeetup implements meetup.MeetupStorage.
func (s *Storage) CancelMeetup(ctx context.Context, meetupID int, cancelledReason string) error {
	query := `
		UPDATE meetup 
		SET 
			status = :status, 
			cancelled_reason = :cancelled_reason, 
			cancelled_at = :cancelled_at, 
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err := s.sqlClient.NamedExecContext(ctx, query, map[string]interface{}{
		"id":               meetupID,
		"status":           entity.StatusCancelled,
		"cancelled_reason": cancelledReason,
		"cancelled_at":     time.Now().Unix(),
		"updated_at":       time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}

	return err
}

// CountOverlappingMeetups implements meetup.MeetupStorage.
func (s *Storage) CountOverlappingMeetups(ctx context.Context, userID int, startTs, endTs int64) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM meetup_user mu
		JOIN meetup m ON mu.meetup_id = m.id
		WHERE mu.user_id = ?
		AND m.start_ts < ? AND m.end_ts > ?
	`

	if err := s.sqlClient.GetContext(ctx, &count, query, userID, endTs, startTs); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("unable to find meetup with user_id %d, start_ts %d, end_ts %d, : %v", userID, startTs, endTs, err)
	}

	return count, nil
}
