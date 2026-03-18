package postgres

import (
	"tech-challenge-user-validation/internal/core/domain"
	"time"
)

type UserModel struct {
	ID        int64     `db:"id"`
	Document  string    `db:"document"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
}

func (UserModel) TableName() string {
	return "users"
}

func (u UserModel) ToDomain() domain.User {
	return domain.User{
		ID:        u.ID,
		Document:  u.Document,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
