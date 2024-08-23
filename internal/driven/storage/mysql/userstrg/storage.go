package userstrg

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

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user UserRow
	query := `
		SELECT
			id,
			username,
			email,
			password,
			created_at
		FROM user
		WHERE username = ?
	`

	if err := s.sqlClient.GetContext(ctx, &user, query, username); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to find user with username %s: %v", username, err)
	}

	return user.ToUser(), nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	var user UserRow
	query := `
		SELECT
			id,
			username,
			email,
			password,
			created_at
		FROM user
		WHERE id = ?
	`

	if err := s.sqlClient.GetContext(ctx, &user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to find user with id %d: %v", id, err)
	}

	return user.ToUser(), nil
}

// JoinMeetup implements meetup.UserStorage.
func (s *Storage) JoinMeetup(ctx context.Context, meetupUser entity.MeetupUser) error {
	meetupUserRow := shared.NewMeetupUserRow(&meetupUser)
	query := `
		REPLACE INTO meetup_user (
			meetup_id, user_id, joined_at
		) VALUES (
			:meetup_id, :user_id, :joined_at
		)
	`

	_, err := s.sqlClient.NamedExecContext(ctx, query, map[string]interface{}{
		"meetup_id": meetupUserRow.MeetupID,
		"user_id":   meetupUserRow.UserID,
		"joined_at": meetupUserRow.JoinedAt,
	})
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}

	return err
}

// CountMeetupUser implements meetup.UserStorage.
func (s *Storage) CountMeetupUser(ctx context.Context, meetupID int, userID int) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM meetup_user
		WHERE meetup_id = ? AND user_id = ?
	`

	if err := s.sqlClient.GetContext(ctx, &count, query, meetupID, userID); err != nil {
		return 0, fmt.Errorf("unable to find supported event with meetup id %d and user id %d: %v", meetupID, userID, err)
	}

	return count, nil
}

// LeaveMeetup implements meetup.UserStorage.
func (s *Storage) LeaveMeetup(ctx context.Context, meetupID, userID int) error {
	query := `
		DELETE FROM meetup_user
		WHERE meetup_id = ? AND user_id = ?
	`

	_, err := s.sqlClient.ExecContext(ctx, query, meetupID, userID)
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}

	return err
}
