package userstrg

import "github.com/Haraj-backend/hex-monscape/internal/core/entity"

type UserRow struct {
	ID        int    `db:"id"`
	Username  string `db:"username"`
	Email     string `db:"email"`
	Password  string `db:"password"`
	CreatedAt int64  `db:"created_at"`
}

func (r *UserRow) ToUser() *entity.User {
	return &entity.User{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}
