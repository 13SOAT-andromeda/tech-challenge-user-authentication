package domain

import "time"

type User struct {
	ID        uint       `json:"id"`
	Password  *Password  `json:"-"`
	Role      string     `json:"role"`
	PersonID  uint       `json:"person_id"`
	Person    *Person    `json:"person,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
