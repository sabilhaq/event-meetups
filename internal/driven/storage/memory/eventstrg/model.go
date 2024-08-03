package eventstrg

import "github.com/Haraj-backend/hex-monscape/internal/core/entity"

type eventRow struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (r eventRow) toEvent() entity.Event {
	return entity.Event{
		ID:   r.ID,
		Name: r.Name,
	}
}
