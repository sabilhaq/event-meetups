package eventstrg

import "github.com/Haraj-backend/hex-monscape/internal/core/entity"

type EventRow struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (r *EventRow) ToEvent() *entity.Event {
	return &entity.Event{
		ID:   r.ID,
		Name: r.Name,
	}
}
