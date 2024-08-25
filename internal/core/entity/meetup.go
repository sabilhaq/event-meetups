package entity

import (
	"gopkg.in/validator.v2"
)

const (
	StatusCancelled = "cancelled"
)

type MeetupConfig struct {
	Name        string `validate:"nonzero"`
	VenueID     int    `validate:"nonzero"`
	EventID     int    `validate:"nonzero"`
	StartTs     int64  `validate:"nonzero"`
	EndTs       int64  `validate:"nonzero"`
	MaxPersons  int    `validate:"nonzero"`
	OrganizerID int    `validate:"nonzero"`
}

func (c MeetupConfig) Validate() error {
	return validator.Validate(c)
}

type CreateMeetupRequest struct {
	Name        string
	VenueID     int
	EventID     int
	StartTs     int64
	EndTs       int64
	MaxPersons  int
	OrganizerID int
}

type Meetup struct {
	ID                 int
	Name               string
	Venue              MeetupVenue
	Event              MeetupEvent
	StartTs            int64
	EndTs              int64
	MaxPersons         int
	Organizer          MeetupOrganizer
	JoinedPersons      []JoinedPerson `json:"joined_persons,omitempty"`
	JoinedPersonsCount int
	IsJoined           bool
	Status             string
	CancelledReason    *string `json:"cancelled_reason,omitempty"`
	CancelledAt        *int64  `json:"cancelled_at,omitempty"`
}

type MeetupVenue struct {
	ID   int
	Name string
}

type MeetupEvent struct {
	ID   int
	Name string
}

type MeetupOrganizer struct {
	ID       int
	Username string
	Email    string
}

type JoinedPerson struct {
	ID       int
	Username string
	Email    string
	JoinedAt int64
}

type GetMeetupFilter struct {
	EventID *int
	Limit   *int
}

type GetMeetupsResponse struct {
	ID                 int
	Name               string
	Venue              MeetupVenue
	Event              MeetupEvent
	StartTs            int64
	EndTs              int64
	MaxPersons         int
	Organizer          MeetupOrganizer
	JoinedPersonsCount int
	Status             string
}

type UpdateMeetupRequest struct {
	Name       string
	StartTs    int64
	EndTs      int64
	MaxPersons int
	UserID     int
}

type CancelMeetupResponse struct {
	ID              int
	Name            string
	Venue           MeetupVenue
	Event           MeetupEvent
	StartTs         int64
	EndTs           int64
	MaxPersons      int
	Organizer       MeetupOrganizer
	Status          string
	CancelledReason string
	CancelledAt     int64
}

func NewMeetup(cfg MeetupConfig) (*Meetup, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	m := &Meetup{
		Name:          cfg.Name,
		Venue:         MeetupVenue{ID: cfg.VenueID},
		Event:         MeetupEvent{ID: cfg.EventID},
		StartTs:       cfg.StartTs,
		EndTs:         cfg.EndTs,
		MaxPersons:    cfg.MaxPersons,
		IsJoined:      false,
		Organizer:     MeetupOrganizer{ID: cfg.OrganizerID},
		JoinedPersons: []JoinedPerson{},
		Status:        "open",
	}
	return m, nil
}

type MeetupUser struct {
	MeetupID int
	UserID   int
	JoinedAt int64
}

type GetIncomingMeetupFilter struct {
	UserID   int
	EventIDs *string
	VenueIDs *string
	Status   string
}
