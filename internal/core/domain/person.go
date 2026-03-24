package domain

import "time"

type Person struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Contact   string     `json:"contact"`
	Document  string     `json:"document"`
	IsActive  bool       `json:"is_active"`
	Address   *Address   `json:"address,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
