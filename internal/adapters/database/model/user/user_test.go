package user

import (
	"testing"
	"time"

	"tech-challenge-user-validation/internal/adapters/database/model/address"
	"tech-challenge-user-validation/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserModel_ToDomain(t *testing.T) {
	now := time.Now()
	model := &Model{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:     "John Doe",
		Email:    "john@example.com",
		Contact:  "123456789",
		Document: "123.456.789-00",
		IsActive: true,
		Password: "hashed_password",
		Role:     "user",
		Address: &address.Model{
			Street: "Main St",
		},
	}

	domainUser := model.ToDomain()

	assert.NotNil(t, domainUser)
	assert.Equal(t, uint(1), domainUser.ID)
	assert.Equal(t, "John Doe", domainUser.Name)
	assert.Equal(t, "john@example.com", domainUser.Email)
	assert.Equal(t, "123456789", domainUser.Contact)
	assert.Equal(t, "123.456.789-00", domainUser.Document)
	assert.True(t, domainUser.IsActive)
	assert.Equal(t, "user", domainUser.Role)
	assert.Equal(t, "hashed_password", domainUser.Password.GetHashed())
	assert.NotNil(t, domainUser.Address)
	assert.Equal(t, "Main St", domainUser.Address.Street)
	assert.Equal(t, now, domainUser.CreatedAt)
	assert.Nil(t, domainUser.DeletedAt)
}

func TestUserModel_FromDomain(t *testing.T) {
	now := time.Now()
	domainUser := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Contact:   "123456789",
		Document:  "123.456.789-00",
		IsActive:  true,
		Role:      "user",
		Address: &domain.Address{
			Street: "Main St",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	model := &Model{}
	model.FromDomain(domainUser) // Note: Password hashing in FromDomain will hash an empty string if not set

	assert.Equal(t, uint(1), model.ID)
	assert.Equal(t, "John Doe", model.Name)
	assert.Equal(t, "john@example.com", model.Email)
	assert.Equal(t, "123456789", model.Contact)
	assert.Equal(t, "123.456.789-00", model.Document)
	assert.True(t, model.IsActive)
	assert.Equal(t, "user", model.Role)
	assert.NotNil(t, model.Address)
	assert.Equal(t, "Main St", model.Address.Street)
	assert.Equal(t, now, model.CreatedAt)
	assert.False(t, model.DeletedAt.Valid)
}

func TestUserModel_FromDomain_Nil(t *testing.T) {
	model := &Model{Name: "Old Name"}
	model.FromDomain(nil)
	assert.Equal(t, "Old Name", model.Name)
}
