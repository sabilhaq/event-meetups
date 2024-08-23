package meetup

import (
	"context"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type MeetupStorage interface {
	// CountMeetups returns existing meetups count for given venueID, eventID, startTs, and endTs from storage. Returns zero
	// when given filter is not found in database.
	CountMeetups(ctx context.Context, venueID, eventID int, startTs, endTs int64) (*int, error)

	// SaveMeetup is used for save meetup in storage.
	SaveMeetup(ctx context.Context, meetup entity.Meetup) (int, error)

	// GetMeetups returns list of meetup available in the system.
	// Returns nil when there is no meetups available.
	GetMeetups(ctx context.Context, filter entity.GetMeetupFilter) ([]entity.Meetup, error)

	// GetMeetup returns meetup instance for given meetupID from storage. Returns nil
	// when given meetupID is not found in database.
	GetMeetup(ctx context.Context, meetupID, userID int) (*entity.Meetup, bool, error)

	// CancelMeetup is used to update meetup status to cancelled in storage.
	CancelMeetup(ctx context.Context, meetupID int, cancelledReason string) error

	// CountOverlappingMeetups is used to count overlapping meetup given userID, startTs, and endTs from storage.
	CountOverlappingMeetups(ctx context.Context, userID int, startTs, endTs int64) (int, error)
}

type VenueStorage interface {
	// IsEventSupported returns true if event supported by the venue.
	// Returns false otherwise.
	IsEventSupported(ctx context.Context, venueID, eventID int) (bool, error)

	// GetVenueCapacity returns venue capacity for given venueID and eventID from storage. Returns zero
	// when given venueID is not found in database.
	GetVenueCapacity(ctx context.Context, venueID, eventID int) (*int, error)

	// GetVenue returns venue instance for given venueID from storage. Returns nil
	// when given venueID is not found in database.
	GetVenue(ctx context.Context, venueID int) (*entity.Venue, error)
}

type EventStorage interface {
	// GetEvent returns event instance for given eventID from storage. Returns nil
	// when given eventID is not found in database.
	GetEvent(ctx context.Context, eventID int) (*entity.Event, error)
}

type UserStorage interface {
	// GetUserByID returns user instance for given userID from storage. Returns nil
	// when given userID is not found in database.
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)

	// JoinMeetup insert a row to table meetup_user
	JoinMeetup(ctx context.Context, meetupUser entity.MeetupUser) error

	// LeaveMeetup delete a row from table meetup_user
	LeaveMeetup(ctx context.Context, meetupID, userID int) error

	// CountMeetupUser checks is user joined to a meetup
	CountMeetupUser(ctx context.Context, meetupID, userID int) (int, error)
}
