package rest

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int
	Err        string
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v - %v - %v", e.StatusCode, e.Err, e.Message)
}

func (e *Error) Is(target error) bool {
	var restErr *Error
	if !errors.As(target, &restErr) {
		return false
	}
	return *e == *restErr
}

func NewInternalServerError(msg string) *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		Err:        "ERR_INTERNAL_ERROR",
		Message:    msg,
	}
}

func NewBadRequestError(msg string) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Err:        "ERR_BAD_REQUEST",
		Message:    fmt.Sprintf("invalid value of `%s`", msg),
	}
}

func NewUnauthorizedError() *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,
		Err:        "ERR_INVALID_ACCESS_TOKEN",
		Message:    "invalid access token",
	}
}

func NewForbiddenError() *Error {
	return &Error{
		StatusCode: http.StatusForbidden,
		Err:        "ERR_FORBIDDEN_ACCESS",
		Message:    "user doesn't have enough authorization",
	}
}

func NewNotFoundError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_NOT_FOUND",
		Message:    "resource is not found",
	}
}

func NewPartnerNotFoundError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_PARTNER_NOT_FOUND",
		Message:    "given `partner_id` is not found",
	}
}

func NewGameNotFoundError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_GAME_NOT_FOUND",
		Message:    "game is not found",
	}
}

func NewBattleNotFoundError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_BATTLE_NOT_FOUND",
		Message:    "battle is not found",
	}
}

func NewInvalidBattleStateError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_INVALID_BATTLE_STATE",
		Message:    "invalid battle state",
	}
}

func NewSessionInvalidCredsError() *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Err:        "ERR_INVALID_CREDS",
		Message:    "invalid username or password",
	}
}

func NewInvalidEventError() *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Err:        "ERR_INVALID_EVENT",
		Message:    "event is not supported by the venue",
	}
}

func NewExceedVenueCapacityError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_EXCEED_VENUE_CAPACITY",
		Message:    "venue capacity is full on the designated meetup time",
	}
}

func NewVenueIsClosedError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_VENUE_IS_CLOSED",
		Message:    "venue is closed on the designated meetup time",
	}
}

func NewMaxPersonsLessThanJoinedPersonsError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_MAX_PERSONS_LESS_THAN_JOINED_PERSONS",
		Message:    "max persons is less than number of joined persons",
	}
}

func NewCancelledReasonRequiredError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_CANCELLED_REASON_REQUIRED",
		Message:    "Cancelled reason is required",
	}
}

func NewMeetupStartedError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_MEETUP_STARTED",
		Message:    "Meetup is started",
	}
}

func NewMeetupFinishedError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_MEETUP_FINISHED",
		Message:    "Meetup is finished",
	}
}

func NewMeetupCancelledError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_MEETUP_CANCELLED",
		Message:    "Meetup is cancelled",
	}
}

func NewMeetupClosedError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_MEETUP_CLOSED",
		Message:    "Meetup is closed",
	}
}

func NewMeetupOverlapsError() *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Err:        "ERR_MEETUP_OVERLAPS",
		Message:    "Meetup overlaps with other meetup that user already joined",
	}
}

func NewUserNotParticipantError() *Error {
	return &Error{
		StatusCode: http.StatusForbidden,
		Err:        "ERR_USER_NOT_PARTICIPANT",
		Message:    "User is not a participant",
	}
}
