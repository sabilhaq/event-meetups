package meetupstrg

import (
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type MeetupRow struct {
	ID              int    `db:"id"`
	Name            string `db:"name"`
	VenueID         int    `db:"venue_id"`
	EventID         int    `db:"event_id"`
	StartTs         int64  `db:"start_ts"`
	EndTs           int64  `db:"end_ts"`
	MaxPersons      int    `db:"max_persons"`
	OrganizerID     int    `db:"organizer_id"`
	Status          string `db:"status"`
	CancelledReason string `db:"cancelled_reason"`
	CancelledAt     *int64 `db:"cancelled_at"`
	CreatedAt       int64  `db:"created_at"`
	UpdatedAt       *int64 `db:"updated_at"`
}

type MeetupJoinVenueEventRow struct {
	ID                int    `db:"id"`
	Name              string `db:"name"`
	VenueID           int    `db:"venue_id"`
	VenueName         string `db:"venue_name"`
	EventID           int    `db:"event_id"`
	EventName         string `db:"event_name"`
	StartTs           int64  `db:"start_ts"`
	EndTs             int64  `db:"end_ts"`
	MaxPersons        int    `db:"max_persons"`
	OrganizerID       int    `db:"organizer_id"`
	OrganizerUsername string `db:"organizer_username"`
	OrganizerEmail    string `db:"organizer_email"`
	Status            string `db:"status"`
	CancelledReason   string `db:"cancelled_reason"`
	CancelledAt       *int64 `db:"cancelled_at"`
	CreatedAt         int64  `db:"created_at"`
	UpdatedAt         int64  `db:"updated_at"`
}

type MeetupJoinVenueEventRows []MeetupJoinVenueEventRow

func (r *MeetupJoinVenueEventRow) ToMeetup() *entity.Meetup {
	return &entity.Meetup{
		ID:         r.ID,
		Name:       r.Name,
		Venue:      entity.MeetupVenue{ID: r.VenueID, Name: r.VenueName},
		Event:      entity.MeetupEvent{ID: r.EventID, Name: r.EventName},
		StartTs:    r.StartTs,
		EndTs:      r.EndTs,
		MaxPersons: r.MaxPersons,
		Organizer:  entity.MeetupOrganizer{ID: r.OrganizerID, Username: r.OrganizerUsername, Email: r.OrganizerEmail},
		Status:     r.Status,
	}
}

func (r MeetupJoinVenueEventRows) ToMeetups() []entity.Meetup {
	var meetups []entity.Meetup
	for _, row := range r {
		meetups = append(meetups, *row.ToMeetup())
	}
	return meetups
}

func NewMeetupRow(meetup *entity.Meetup) *MeetupRow {
	now := time.Now().Unix()

	return &MeetupRow{
		ID:              meetup.ID,
		Name:            meetup.Name,
		VenueID:         meetup.Venue.ID,
		EventID:         meetup.Event.ID,
		StartTs:         meetup.StartTs,
		EndTs:           meetup.EndTs,
		MaxPersons:      meetup.MaxPersons,
		OrganizerID:     meetup.Organizer.ID,
		Status:          meetup.Status,
		CancelledReason: "",
		CancelledAt:     nil,
		CreatedAt:       now,
		UpdatedAt:       &now,
	}
}
