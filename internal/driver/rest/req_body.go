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
