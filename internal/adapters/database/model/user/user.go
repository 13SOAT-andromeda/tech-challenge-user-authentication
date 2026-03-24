package user

import (
	"time"

	personModel "tech-challenge-user-validation/internal/adapters/database/model/person"
	"tech-challenge-user-validation/internal/core/domain"
	"tech-challenge-user-validation/pkg/encryption"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Password string            `gorm:"not null"`
	Role     string            `gorm:"not null"`
	PersonID uint              `gorm:"not null"`
	Person   personModel.Model `gorm:"foreignKey:PersonID;references:ID"`
}

func (*Model) TableName() string {
	return "User"
}

func (m *Model) ToDomain() *domain.User {
	pass := domain.NewPasswordFromHash(m.Password, encryption.NewBcryptHasher())

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	var person *domain.Person
	if m.Person.ID != 0 {
		person = m.Person.ToDomain()
	}

	return &domain.User{
		ID:        m.ID,
		Password:  &pass,
		Role:      m.Role,
		PersonID:  m.PersonID,
		Person:    person,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (m *Model) FromDomain(d *domain.User) {
	if d == nil {
		return
	}

	m.ID = d.ID
	m.Role = d.Role
	m.PersonID = d.PersonID

	if d.Password != nil {
		m.Password = d.Password.GetHashed()
	}

	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt

	if d.DeletedAt != nil {
		m.DeletedAt = gorm.DeletedAt{Time: *d.DeletedAt, Valid: true}
	} else {
		m.DeletedAt = gorm.DeletedAt{Valid: false}
	}

	if d.Person != nil {
		m.Person.FromDomain(d.Person)
	}
}
