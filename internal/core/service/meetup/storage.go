package meetup

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type MeetupStorage interface {
	// GetMeetups returns list of meetup available in the system.
	// Returns nil when there is no meetups available.
	GetMeetups(ctx context.Context) ([]entity.Meetup, error)

	// SaveMeetup is used for save meetup in storage.
	SaveMeetup(ctx context.Context, meetup entity.Meetup) error

	// GetMeetup returns meetup instance for given meetupID from storage. Returns nil
	// when given meetupID is not found in database.
	GetMeetup(ctx context.Context, meetupID int) (*entity.Meetup, error)

	// CancelMeetup is used to update meetup status to cancelled in storage.
	CancelMeetup(ctx context.Context, meetupID int, cancelledReason string) error
}
