package mailpit

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"gopkg.in/validator.v2"
)

type Storage struct {
	host      string
	port      string
	fromEmail string
}

type Config struct {
	Host      string `validate:"nonnil"`
	Port      string `validate:"nonnil"`
	FromEmail string `validate:"nonnil"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &Storage{
		host:      cfg.Host,
		port:      cfg.Port,
		fromEmail: cfg.FromEmail,
	}
	return s, nil
}

func (s *Storage) SendCancellationEmail(toEmails []string, cancelledReason string) error {
	to := strings.Join(toEmails, ", ")
	subject := "Meetup Cancellation Notice"
	body := fmt.Sprintf("We regret to inform you that the meetup you joined has been canceled for the following reason:\n\n%s", cancelledReason)

	// Construct the email message
	msg := []byte("To: " + to + "\r\n" +
		"From: " + s.fromEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	if err := smtp.SendMail(s.host+":"+s.port, nil, s.fromEmail, toEmails, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

func (s *Storage) NotifyOrganizer(organizerEmail, username string, joinedCount int) error {
	to := []string{organizerEmail}
	subject := "New User Joined Your Meetup"
	body := fmt.Sprintf("User %s just joined your meetup. Current number of joined persons: %d.", username, joinedCount)

	// Construct the email message
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", s.fromEmail, to[0], subject, body)

	if err := smtp.SendMail(s.host+":"+s.port, nil, s.fromEmail, to, []byte(message)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", organizerEmail)
	return nil
}
