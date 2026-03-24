package person

import (
	"time"

	"tech-challenge-user-validation/internal/adapters/database/model/address"
	"tech-challenge-user-validation/internal/core/domain"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name     string         `gorm:"not null"`
	Email    string         `gorm:"not null;unique"`
	Contact  string         `gorm:"not null"`
	Document string         `gorm:"not null;unique"`
	IsActive bool           `gorm:"default:true"`
	Address  *address.Model `gorm:"embedded"`
}

func (*Model) TableName() string {
	return "Person"
}

func (m *Model) ToDomain() *domain.Person {
	var addressDomain *domain.Address
	if m.Address != nil {
		addressDomain = m.Address.ToDomain()
	}
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}
	return &domain.Person{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Contact:   m.Contact,
		Document:  m.Document,
		IsActive:  m.IsActive,
		Address:   addressDomain,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (m *Model) FromDomain(d *domain.Person) {
	if d == nil {
		return
	}
	m.ID = d.ID
	m.Name = d.Name
	m.Email = d.Email
	m.Contact = d.Contact
	m.Document = d.Document
	m.IsActive = d.IsActive
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	if d.DeletedAt != nil {
		m.DeletedAt = gorm.DeletedAt{Time: *d.DeletedAt, Valid: true}
	} else {
		m.DeletedAt = gorm.DeletedAt{Valid: false}
	}
	if m.Address == nil {
		m.Address = &address.Model{}
	}
	m.Address.FromDomain(d.Address)
}
