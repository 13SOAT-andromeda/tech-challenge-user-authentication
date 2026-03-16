package domain

import "time"

type User struct {
	ID        int64     `json:"id"`
	CPF       string    `json:"cpf"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
