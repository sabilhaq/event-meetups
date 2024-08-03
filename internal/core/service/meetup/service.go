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
	ErrMeetupNotFound = errors.New("meetup is not found")
)

type Service interface {
	// CreateMeetup is used to add a new meetup to the system. If the given `meetupID` not found in storage,
	// it returns `ErrMeetupNotFound`. Upon success it returns meetup instance that being saved on storage.
	// CreateMeetup is used to add a new meetup to the system. Meetup is a gathering event in a venue in a specific range of time.
	// Meetup can only be created in a venue that supports the event, within the operating hours of the venue, and not exceeding the capacity of the venue in that time.
	CreateMeetup(ctx context.Context, req entity.CreateMeetupRequest) (*entity.Meetup, error)

	// GetMeetups returns all meetups available in the system.
	GetMeetups(ctx context.Context) ([]entity.GetMeetupsResponse, error)

	// GetMeetup returns a single meetup from storage from given meetup id. Upon meetup is not found, it returns
	// `ErrMeetupNotFound`.
	GetMeetup(ctx context.Context, meetupID int) (*entity.Meetup, error)

	// UpdateMeetup is used to update a meetup. Only the organizer of the meetup can update the meetup. Update action is limited to:
	// - Change the name of the meetup
	// - Change the start and end time of the meetup, but it still need to follows the rules of Create Meetup endpoint
	// - Update maximum number of persons that can join the meetup
	UpdateMeetup(ctx context.Context, meetupID int, req entity.UpdateMeetupRequest) (*entity.Meetup, error)

	// CancelMeetup is used to cancel a meetup. Only the organizer of the meetup can cancel the meetup.
	// Meetup can only be cancelled if it isn't started yet.
	CancelMeetup(ctx context.Context, meetupID int, cancelledReason string) (*entity.CancelMeetupResponse, error)

	// JoinMeetup is used to join a meetup. User can only join a meetup if the meetup is still open
	// which means the meetup hasn't reached the maximum number of persons, not cancelled, and not finished yet.
	JoinMeetup(ctx context.Context, meetupID int) (*entity.Meetup, error)

	// LeaveMeetup is used to leave a meetup. User can only leave a meetup if he/she already
	// joined the meetup, also the meetup is not cancelled or finished yet.
	LeaveMeetup(ctx context.Context, meetupID int) error

	// GetIncomingMeetups is used to list future meetups that are joined by a user. The returned meetup
	// statuses are either open or cancelled.
	GetIncomingMeetups(ctx context.Context) ([]entity.Meetup, error)
}

type service struct {
	meetupStorage MeetupStorage
}

func (s *service) CreateMeetup(ctx context.Context, req entity.CreateMeetupRequest) (*entity.Meetup, error) {
	// TODO: validation
	// initiate new meetup instance
	cfg := ConvertRequestToConfig(req)
	meetup, err := entity.NewMeetup(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize meetup instance due: %w", err)
	}
	// store the meetup instance on storage
	err = s.meetupStorage.SaveMeetup(ctx, *meetup)
	if err != nil {
		return nil, fmt.Errorf("unable to save meetup instance due: %w", err)
	}
	return meetup, nil
}

// Conversion function
func ConvertRequestToConfig(req entity.CreateMeetupRequest) entity.MeetupConfig {
	return entity.MeetupConfig(req)
}

func (s *service) GetMeetups(ctx context.Context) ([]entity.GetMeetupsResponse, error) {
	meetups, err := s.meetupStorage.GetMeetups(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get available meetups due: %w", err)
	}
	res := make([]entity.GetMeetupsResponse, len(meetups))
	for i := 0; i < len(meetups); i++ {
		meetup := entity.GetMeetupsResponse{
			ID: meetups[i].ID,
			// TODO: fill fields
		}
		res = append(res, meetup)
	}
	return res, nil
}

func (s *service) GetMeetup(ctx context.Context, meetupID int) (*entity.Meetup, error) {
	meetup, err := s.getMeetupInstance(ctx, meetupID)
	if meetup == nil {
		return nil, ErrMeetupNotFound
	}
	return meetup, err
}

// getMeetupInstance returns meetup for given meetup id, if meetup is not found
// will be returned nil.
func (s *service) getMeetupInstance(ctx context.Context, meetupID int) (*entity.Meetup, error) {
	meetup, err := s.meetupStorage.GetMeetup(ctx, meetupID)
	if err != nil {
		return nil, fmt.Errorf("unable to get meetup due: %w", err)
	}
	if meetup == nil {
		return nil, nil
	}
	return meetup, nil
}

// UpdateMeetup implements Service.
func (s *service) UpdateMeetup(ctx context.Context, meetupID int, req entity.UpdateMeetupRequest) (*entity.Meetup, error) {
	panic("unimplemented")
}

func (s *service) CancelMeetup(ctx context.Context, meetupID int, cancelledReason string) (*entity.CancelMeetupResponse, error) {
	// get existing meetup
	meetup, err := s.GetMeetup(ctx, meetupID)
	if err != nil {
		return nil, err
	}

	// delete meetup
	err = s.meetupStorage.CancelMeetup(ctx, meetup.ID, cancelledReason)
	if err != nil {
		return nil, fmt.Errorf("unable to delete meetup due: %w", err)
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
		CancelledReason: cancelledReason,
		CancelledAt:     time.Now().Unix(), // TODO: make relation between meetup and meetup_cancelled_reason
	}, nil
}

// JoinMeetup implements Service.
func (s *service) JoinMeetup(ctx context.Context, meetupID int) (*entity.Meetup, error) {
	panic("unimplemented")
}

// LeaveMeetup implements Service.
func (s *service) LeaveMeetup(ctx context.Context, meetupID int) error {
	panic("unimplemented")
}

func (s *service) GetIncomingMeetups(ctx context.Context) ([]entity.Meetup, error) {
	// TODO: validation
	//

	meetups, err := s.meetupStorage.GetMeetups(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get available meetups due: %w", err)
	}
	return meetups, nil
}

type ServiceConfig struct {
	MeetupStorage MeetupStorage `validate:"nonnil"`
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
	}
	return s, nil
}
