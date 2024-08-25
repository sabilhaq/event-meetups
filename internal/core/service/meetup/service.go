package meetup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"gopkg.in/validator.v2"
)

var (
	ErrMeetupNotFound                  = errors.New("meetup is not found")
	ErrInvalidEvent                    = errors.New("event is not supported by the venue")
	ErrExceedVenueCapacity             = errors.New("venue capacity is full on the designated meetup time")
	ErrVenueIsClosed                   = errors.New("venue is closed on the designated meetup time")
	ErrForbidden                       = errors.New("user is not authorized to access this resource")
	ErrMaxPersonsLessThanJoinedPersons = errors.New("max persons is less than number of joined persons")
	ErrCancelledReasonRequred          = errors.New("Cancelled reason is required")
	ErrMeetupStarted                   = errors.New("Meetup is started")
	ErrMeetupFinished                  = errors.New("Meetup is finished")
	ErrMeetupCancelled                 = errors.New("Meetup is cancelled")
	ErrMeetupClosed                    = errors.New("Meetup is closed")
	ErrMeetupOverlaps                  = errors.New("Meetup overlaps with other meetup that user already joined")
	ErrUserNotParticipant              = errors.New("User is not a participant")
)

type Service interface {
	// CreateMeetup is used to add a new meetup to the system. If the given `meetupID` not found in storage,
	// it returns `ErrMeetupNotFound`. Upon success it returns meetup instance that being saved on storage.
	// CreateMeetup is used to add a new meetup to the system. Meetup is a gathering event in a venue in a specific range of time.
	// Meetup can only be created in a venue that supports the event, within the operating hours of the venue, and not exceeding the capacity of the venue in that time.
	CreateMeetup(ctx context.Context, req entity.CreateMeetupRequest) (*entity.Meetup, error)

	// GetMeetups returns all meetups available in the system.
	// The result is sorted from the nearest meetup time to the furthest.
	GetMeetups(ctx context.Context, filter entity.GetMeetupFilter) ([]entity.GetMeetupsResponse, error)

	// GetMeetup returns a single meetup from storage from given meetup id. Upon meetup is not found, it returns
	// `ErrMeetupNotFound`.
	GetMeetup(ctx context.Context, meetupID, userID int) (*entity.Meetup, error)

	// UpdateMeetup is used to update a meetup. Only the organizer of the meetup can update the meetup. Update action is limited to:
	// - Change the name of the meetup
	// - Change the start and end time of the meetup, but it still need to follows the rules of Create Meetup endpoint
	// - Update maximum number of persons that can join the meetup
	UpdateMeetup(ctx context.Context, meetupID int, req entity.UpdateMeetupRequest) (*entity.Meetup, error)

	// CancelMeetup is used to cancel a meetup. Only the organizer of the meetup can cancel the meetup.
	// Meetup can only be cancelled if it isn't started yet.
	CancelMeetup(ctx context.Context, meetupID int, userID int, cancelledReason string) (*entity.CancelMeetupResponse, error)

	// JoinMeetup is used to join a meetup. User can only join a meetup if the meetup is still open
	// which means the meetup hasn't reached the maximum number of persons, not cancelled, and not finished yet.
	JoinMeetup(ctx context.Context, meetupID int, userID int) (*entity.Meetup, error)

	// LeaveMeetup is used to leave a meetup. User can only leave a meetup if he/she already
	// joined the meetup, also the meetup is not cancelled or finished yet.
	LeaveMeetup(ctx context.Context, meetupID int, userID int) error

	// GetIncomingMeetups is used to list future meetups that are joined by a user. The returned meetup
	// statuses are either open or cancelled.
	GetIncomingMeetups(ctx context.Context, filter entity.GetIncomingMeetupFilter) ([]entity.Meetup, error)
}

type service struct {
	meetupStorage MeetupStorage
	venueStorage  VenueStorage
	eventStorage  EventStorage
	userStorage   UserStorage
	emailStorage  EmailStorage
}

func (s *service) CreateMeetup(ctx context.Context, req entity.CreateMeetupRequest) (*entity.Meetup, error) {
	// check is event supported by the venue
	isEventSupported, err := s.venueStorage.IsEventSupported(ctx, req.VenueID, req.EventID)
	if err != nil {
		return nil, fmt.Errorf("unable to check is event supported by the venue due: %w", err)
	}
	if !isEventSupported {
		return nil, ErrInvalidEvent
	}

	// check is venue capacity exceed
	venueCapacity, err := s.venueStorage.GetVenueCapacity(ctx, req.VenueID, req.EventID)
	if err != nil {
		return nil, fmt.Errorf("unable to get supported venue due: %w", err)
	}
	if venueCapacity == nil {
		return nil, ErrInvalidEvent
	}

	existingMeetupsCount, err := s.meetupStorage.CountMeetups(ctx, req.VenueID, req.EventID, req.StartTs, req.EndTs)
	if err != nil || existingMeetupsCount == nil {
		return nil, fmt.Errorf("unable to get existing meetups count due: %w", err)
	}

	if *existingMeetupsCount >= *venueCapacity {
		return nil, ErrExceedVenueCapacity
	}

	// check is meetup within venue operating hours
	venue, err := s.venueStorage.GetVenue(ctx, req.VenueID)
	if err != nil {
		return nil, fmt.Errorf("unable to get existing meetups count due: %w", err)
	}

	// Convert timestamp to the venue's timezone
	loc, err := time.LoadLocation(venue.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %v", err)
	}
	startTime := time.Unix(req.StartTs, 0).In(loc)
	endTime := time.Unix(req.EndTs, 0).In(loc)

	// get day of the week for start and end time
	startDay := startTime.Weekday()
	endDay := endTime.Weekday()

	// check if the venue is open on the days of the meetup
	openDays := make(map[int]bool)
	for _, day := range venue.OpenDays {
		openDays[day] = true
	}

	if !openDays[int(startDay)] || !openDays[int(endDay)] {
		return nil, ErrVenueIsClosed // meetup is not within venue's operating days
	}

	// convert venue open and close time to full time in venue's timezone
	openTime, err := time.ParseInLocation("15:04", venue.OpenAt, loc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse open time: %v", err)
	}
	closeTime, err := time.ParseInLocation("15:04", venue.ClosedAt, loc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse close time: %v", err)
	}

	// extract the time of day from the meetup start and end times
	startHour := time.Date(0, 1, 1, startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
	endHour := time.Date(0, 1, 1, endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

	// check if meetup start and end times are within the operating hours (ignoring the date)
	if startHour.Before(openTime) || endHour.After(closeTime) {
		return nil, ErrVenueIsClosed // Meetup is not within venue's operating hours
	}

	// initiate new meetup instance
	cfg := entity.MeetupConfig(req)
	meetup, err := entity.NewMeetup(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize meetup instance due: %w", err)
	}
	// store the meetup instance on storage
	id, err := s.meetupStorage.SaveMeetup(ctx, *meetup)
	if err != nil {
		return nil, fmt.Errorf("unable to save meetup instance due: %w", err)
	}

	event, err := s.eventStorage.GetEvent(ctx, req.EventID)
	if err != nil {
		return nil, fmt.Errorf("unable to get event due: %w", err)
	}

	organizer, err := s.userStorage.GetUserByID(ctx, req.OrganizerID)
	if err != nil {
		return nil, fmt.Errorf("unable to get user due: %w", err)
	}

	meetup.ID = id
	meetup.Venue.ID = venue.ID
	meetup.Venue.Name = venue.Name
	meetup.Event.ID = event.ID
	meetup.Event.Name = event.Name
	meetup.Organizer.ID = organizer.ID
	meetup.Organizer.Username = organizer.Username
	meetup.Organizer.Email = organizer.Email

	return meetup, nil
}

func (s *service) GetMeetups(ctx context.Context, filter entity.GetMeetupFilter) ([]entity.GetMeetupsResponse, error) {
	meetups, err := s.meetupStorage.GetMeetups(ctx, entity.GetMeetupFilter{
		EventID: filter.EventID,
		Limit:   filter.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get available meetups due: %w", err)
	}
	res := make([]entity.GetMeetupsResponse, len(meetups))
	for i := 0; i < len(meetups); i++ {
		meetup := entity.GetMeetupsResponse{
			ID:                 meetups[i].ID,
			Name:               meetups[i].Name,
			Venue:              meetups[i].Venue,
			Event:              meetups[i].Event,
			StartTs:            meetups[i].StartTs,
			EndTs:              meetups[i].EndTs,
			MaxPersons:         meetups[i].MaxPersons,
			Organizer:          meetups[i].Organizer,
			JoinedPersonsCount: meetups[i].JoinedPersonsCount,
			Status:             meetups[i].Status,
		}
		res[i] = meetup
	}
	return res, nil
}

// GetMeetup returns meetup for given meetup id, if meetup is not found
// will be returned nil.
func (s *service) GetMeetup(ctx context.Context, meetupID, userID int) (*entity.Meetup, error) {
	meetup, isOrganizerOrParticipant, err := s.meetupStorage.GetMeetup(ctx, meetupID, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get meetup due: %w", err)
	}
	if meetup == nil {
		return nil, ErrMeetupNotFound
	}
	if len(meetup.JoinedPersons) == 0 {
		meetup.JoinedPersons = []entity.JoinedPerson{}
	}
	if !isOrganizerOrParticipant {
		// Clear the JoinedPersons slice and Cancelled fields if the user isn't an organizer or participant
		meetup.JoinedPersons = nil
		meetup.CancelledReason = nil
		meetup.CancelledAt = nil
	}
	return meetup, err
}

// UpdateMeetup implements Service.
func (s *service) UpdateMeetup(ctx context.Context, meetupID int, req entity.UpdateMeetupRequest) (*entity.Meetup, error) {
	meetup, _, err := s.meetupStorage.GetMeetup(ctx, meetupID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("unable to get meetup due: %w", err)
	}
	if meetup == nil {
		return nil, ErrMeetupNotFound
	}
	if meetup.Organizer.ID != req.UserID {
		return nil, ErrForbidden
	}
	if req.MaxPersons < meetup.JoinedPersonsCount {
		return nil, ErrMaxPersonsLessThanJoinedPersons
	}

	// check is event supported by the venue
	isEventSupported, err := s.venueStorage.IsEventSupported(ctx, meetup.Venue.ID, meetup.Event.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to check is event supported by the venue due: %w", err)
	}
	if !isEventSupported {
		return nil, ErrInvalidEvent
	}

	// check is venue capacity exceed
	venueCapacity, err := s.venueStorage.GetVenueCapacity(ctx, meetup.Venue.ID, meetup.Event.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get supported venue due: %w", err)
	}
	if venueCapacity == nil {
		return nil, ErrInvalidEvent
	}

	existingMeetupsCount, err := s.meetupStorage.CountMeetups(ctx, meetup.Venue.ID, meetup.Event.ID, req.StartTs, req.EndTs)
	if err != nil || existingMeetupsCount == nil {
		return nil, fmt.Errorf("unable to get existing meetups count due: %w", err)
	}

	if *existingMeetupsCount >= *venueCapacity {
		return nil, ErrExceedVenueCapacity
	}

	// check is meetup within venue operating hours
	venue, err := s.venueStorage.GetVenue(ctx, meetup.Venue.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get existing meetups count due: %w", err)
	}

	// Convert timestamp to the venue's timezone
	loc, err := time.LoadLocation(venue.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %v", err)
	}
	startTime := time.Unix(req.StartTs, 0).In(loc)
	endTime := time.Unix(req.EndTs, 0).In(loc)

	// get day of the week for start and end time
	startDay := startTime.Weekday()
	endDay := endTime.Weekday()

	// check if the venue is open on the days of the meetup
	openDays := make(map[int]bool)
	for _, day := range venue.OpenDays {
		openDays[day] = true
	}

	if !openDays[int(startDay)] || !openDays[int(endDay)] {
		return nil, ErrVenueIsClosed // meetup is not within venue's operating days
	}

	// convert venue open and close time to full time in venue's timezone
	openTime, err := time.ParseInLocation("15:04", venue.OpenAt, loc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse open time: %v", err)
	}
	closeTime, err := time.ParseInLocation("15:04", venue.ClosedAt, loc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse close time: %v", err)
	}

	// extract the time of day from the meetup start and end times
	startHour := time.Date(0, 1, 1, startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
	endHour := time.Date(0, 1, 1, endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

	// check if meetup start and end times are within the operating hours (ignoring the date)
	if startHour.Before(openTime) || endHour.After(closeTime) {
		return nil, ErrVenueIsClosed // Meetup is not within venue's operating hours
	}

	id, err := s.meetupStorage.SaveMeetup(ctx, *meetup)
	if err != nil {
		return nil, fmt.Errorf("unable to save meetup instance due: %w", err)
	}

	organizer, err := s.userStorage.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("unable to get user due: %w", err)
	}

	meetup.ID = id
	meetup.Organizer.ID = organizer.ID
	meetup.Organizer.Username = organizer.Username
	meetup.Organizer.Email = organizer.Email

	return meetup, nil
}

func (s *service) CancelMeetup(ctx context.Context, meetupID int, userID int, cancelledReason string) (*entity.CancelMeetupResponse, error) {
	// get existing meetup
	meetup, _, err := s.meetupStorage.GetMeetup(ctx, meetupID, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get meetup due: %w", err)
	}
	if meetup == nil {
		return nil, ErrMeetupNotFound
	}
	if meetup.Organizer.ID != userID {
		return nil, ErrForbidden
	}

	// check is meetup within venue operating hours
	venue, err := s.venueStorage.GetVenue(ctx, meetup.Venue.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get existing meetups count due: %w", err)
	}

	// Convert timestamp to the venue's timezone
	loc, err := time.LoadLocation(venue.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %v", err)
	}
	startTime := time.Unix(meetup.StartTs, 0).In(loc)
	if time.Now().After(startTime) || time.Now().Equal(startTime) {
		return nil, ErrMeetupStarted
	}

	// update meetup status to cancelled
	err = s.meetupStorage.CancelMeetup(ctx, meetup.ID, cancelledReason)
	if err != nil {
		return nil, fmt.Errorf("unable to cancel meetup due: %w", err)
	}

	meetup, _, err = s.meetupStorage.GetMeetup(ctx, meetupID, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get meetup due: %w", err)
	}

	// Get emails of all joined persons
	joinedPersons, err := s.meetupStorage.GetJoinedPersons(ctx, meetup.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get joined persons due: %w", err)
	}

	var emails []string
	for _, user := range joinedPersons {
		emails = append(emails, user.Email)
	}

	// Send cancellation email
	if len(emails) > 0 {
		if err := s.emailStorage.SendCancellationEmail(emails, cancelledReason); err != nil {
			return nil, fmt.Errorf("unable to send emails due: %w", err)
		}
	}

	return &entity.CancelMeetupResponse{
		ID:              meetup.ID,
		Name:            meetup.Name,
		Venue:           meetup.Venue,
		Event:           meetup.Event,
		StartTs:         meetup.StartTs,
		EndTs:           meetup.EndTs,
		MaxPersons:      meetup.MaxPersons,
		Organizer:       meetup.Organizer,
		Status:          meetup.Status,
		CancelledReason: *meetup.CancelledReason,
		CancelledAt:     *meetup.CancelledAt,
	}, nil
}

// JoinMeetup implements Service.
func (s *service) JoinMeetup(ctx context.Context, meetupID int, userID int) (*entity.Meetup, error) {
	meetup, _, err := s.meetupStorage.GetMeetup(ctx, meetupID, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get meetup due: %w", err)
	}
	if meetup == nil {
		return nil, ErrMeetupNotFound
	}

	venue, err := s.venueStorage.GetVenue(ctx, meetup.Venue.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get existing meetups count due: %w", err)
	}
	loc, err := time.LoadLocation(venue.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %v", err)
	}
	endTime := time.Unix(meetup.EndTs, 0).In(loc)
	if time.Now().After(endTime) || time.Now().Equal(endTime) {
		return nil, ErrMeetupFinished
	}

	if meetup.Status == entity.StatusCancelled {
		return nil, ErrMeetupCancelled
	}

	if meetup.JoinedPersonsCount == meetup.MaxPersons {
		return nil, ErrMeetupClosed
	}

	overlapCount, err := s.meetupStorage.CountOverlappingMeetups(ctx, userID, meetup.StartTs, meetup.EndTs)
	if err != nil {
		return nil, err
	}
	if overlapCount > 0 {
		return nil, ErrMeetupOverlaps
	}

	err = s.userStorage.JoinMeetup(ctx, entity.MeetupUser{
		MeetupID: meetup.ID,
		UserID:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to join meetup due: %w", err)
	}

	meetup, _, err = s.meetupStorage.GetMeetup(ctx, meetupID, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get meetup due: %w", err)
	}

	user, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get user due: %w", err)
	}

	// Send notification to the meetup organizer
	if err := s.emailStorage.NotifyOrganizer(meetup.Organizer.Email, user.Username, len(meetup.JoinedPersons)); err != nil {
		return nil, fmt.Errorf("unable to send emails due: %w", err)
	}

	return meetup, nil
}

// LeaveMeetup implements Service.
func (s *service) LeaveMeetup(ctx context.Context, meetupID int, userID int) error {
	meetup, _, err := s.meetupStorage.GetMeetup(ctx, meetupID, userID)
	if err != nil {
		return fmt.Errorf("unable to get meetup due: %w", err)
	}
	if meetup == nil {
		return ErrMeetupNotFound
	}

	venue, err := s.venueStorage.GetVenue(ctx, meetup.Venue.ID)
	if err != nil {
		return fmt.Errorf("unable to get existing meetups count due: %w", err)
	}
	loc, err := time.LoadLocation(venue.Timezone)
	if err != nil {
		return fmt.Errorf("failed to load location: %v", err)
	}
	endTime := time.Unix(meetup.EndTs, 0).In(loc)
	if time.Now().After(endTime) || time.Now().Equal(endTime) {
		return ErrMeetupFinished
	}

	if meetup.Status == entity.StatusCancelled {
		return ErrMeetupCancelled
	}

	count, err := s.userStorage.CountMeetupUser(ctx, meetup.ID, userID)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrUserNotParticipant
	}

	err = s.userStorage.LeaveMeetup(ctx, meetup.ID, userID)
	if err != nil {
		return fmt.Errorf("unable to leave meetup due: %w", err)
	}

	return nil
}

func (s *service) GetIncomingMeetups(ctx context.Context, filter entity.GetIncomingMeetupFilter) ([]entity.Meetup, error) {
	meetups, err := s.meetupStorage.GetIncomingMeetups(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to get incoming meetups due: %w", err)
	}
	return meetups, nil
}

type ServiceConfig struct {
	MeetupStorage MeetupStorage `validate:"nonnil"`
	VenueStorage  VenueStorage  `validate:"nonnil"`
	EventStorage  EventStorage  `validate:"nonnil"`
	UserStorage   UserStorage   `validate:"nonnil"`
	EmailStorage  EmailStorage  `validate:"nonnil"`
}

func (c ServiceConfig) Validate() error {
	return validator.Validate(c)
}

// NewService returns new instance of service.
func NewService(cfg ServiceConfig) (Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &service{
		meetupStorage: cfg.MeetupStorage,
		venueStorage:  cfg.VenueStorage,
		eventStorage:  cfg.EventStorage,
		userStorage:   cfg.UserStorage,
		emailStorage:  cfg.EmailStorage,
	}
	return s, nil
}
