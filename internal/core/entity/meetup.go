package entity

import (
	"gopkg.in/validator.v2"
)

type MeetupConfig struct {
	Name       string `validate:"nonzero"`
	VenueID    int    `validate:"nonzero"`
	EventID    int    `validate:"nonzero"`
	StartTs    int64  `validate:"nonzero"`
	EndTs      int64  `validate:"nonzero"`
	MaxPersons int    `validate:"nonzero"`
}

func (c MeetupConfig) Validate() error {
	return validator.Validate(c)
}

type CreateMeetupRequest struct {
	Name       string
	VenueID    int
	EventID    int
	StartTs    int64
	EndTs      int64
	MaxPersons int
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
	JoinedPersons      []JoinedPerson
	JoinedPersonsCount int
	IsJoined           bool
	Status             string
	CancelledReason    string
	CancelledAt        *int64
	CreatedAt          int64
	UpdatedAt          *int64
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
	ID       string
	Username string
	Email    string
	JoinedAt int
}

type GetMeetupsResponse struct {
	ID                 int
	Name               string
	Venue              MeetupVenue
	Event              MeetupEvent
	StartTs            int
	EndTs              int
	MaxPersons         int
	Organizer          MeetupOrganizer
	JoinedPersonsCount int
	Status             string
}

type UpdateMeetupRequest struct {
	ID         string
	Name       string
	StartTs    int
	EndTs      int
	MaxPersons int
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
		Name:       cfg.Name,
		Venue:      MeetupVenue{ID: cfg.VenueID},
		Event:      MeetupEvent{ID: cfg.EventID},
		StartTs:    cfg.StartTs,
		EndTs:      cfg.EndTs,
		MaxPersons: cfg.MaxPersons,
		IsJoined:   false,
		Status:     "open",
	}
	return m, nil
}
