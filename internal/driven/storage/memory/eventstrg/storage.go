package eventstrg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"gopkg.in/validator.v2"
)

type Storage struct {
	data map[int]entity.Event
}

// GetEvents implements event.EventStorage.
func (s *Storage) GetEvents(ctx context.Context) ([]entity.Event, error) {
	var events []entity.Event
	for _, event := range s.data {
		events = append(events, event)
	}
	return events, nil
}

type Config struct {
	EventData []byte `validate:"nonzero"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	// validate config
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	// parse event data
	var rows []eventRow
	err = json.Unmarshal(cfg.EventData, &rows)
	if err != nil {
		return nil, fmt.Errorf("unable to parse event data due: %w", err)
	}
	data := map[int]entity.Event{}
	for _, eventRow := range rows {
		event := eventRow.toEvent()
		data[event.ID] = event
	}
	return &Storage{data: data}, nil
}
