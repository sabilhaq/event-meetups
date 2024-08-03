package entity

import "gopkg.in/validator.v2"

type SessionConfig struct {
	UserID   int    `validate:"nonzero"`
	Username string `validate:"nonzero"`
	Email    string `validate:"nonzero"`
	Exp      int64  `validate:"nonzero"`
}

func (c SessionConfig) Validate() error {
	return validator.Validate(c)
}

type Session struct {
	ID          int
	Username    string
	Email       string
	AccessToken string
}

func NewSession(cfg SessionConfig) (*Session, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &Session{
		ID:       cfg.UserID,
		Username: cfg.Username,
		Email:    cfg.Email,
	}
	return s, nil
}
