package rest

import "gopkg.in/validator.v2"

type newGameReqBody struct {
	PlayerName string `json:"player_name" validate:"nonzero"`
	PartnerID  string `json:"partner_id" validate:"nonzero"`
}

func (rb newGameReqBody) Validate() error {
	err := validator.Validate(rb)
	if err != nil {
		return NewBadRequestError(err.Error())
	}
	return nil
}

type newSessionReqBody struct {
	Username string `json:"username" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

func (rb newSessionReqBody) Validate() error {
	err := validator.Validate(rb)
	if err != nil {
		return NewBadRequestError(err.Error())
	}
	return nil
}

type createMeetupReqBody struct {
	Name       string `json:"name" validate:"nonzero"`
	VenueID    int    `json:"venue_id" validate:"nonzero"`
	EventID    int    `json:"event_id" validate:"nonzero"`
	StartTs    int64  `json:"start_ts" validate:"nonzero"`
	EndTs      int64  `json:"end_ts" validate:"nonzero"`
	MaxPersons int    `json:"max_persons" validate:"nonzero"`
}

func (rb createMeetupReqBody) Validate() error {
	err := validator.Validate(rb)
	if err != nil {
		return NewBadRequestError(err.Error())
	}
	return nil
}
