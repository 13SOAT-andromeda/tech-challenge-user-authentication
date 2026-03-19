package domain

import "time"

type User struct {
	ID        uint
	Name      string
	Email     string
	Contact   string
	Role      string
	Document  string
	IsActive  bool
	Password  Password
	Address   *Address
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
