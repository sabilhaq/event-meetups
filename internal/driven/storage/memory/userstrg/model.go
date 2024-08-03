package userstrg

import "github.com/Haraj-backend/hex-monscape/internal/core/entity"

type userRow struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r userRow) toUser() entity.User {
	return entity.User{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}
