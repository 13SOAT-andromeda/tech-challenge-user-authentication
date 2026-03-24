package user

import (
	"testing"
	"time"

	personModel "tech-challenge-user-validation/internal/adapters/database/model/person"
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
		Password: "hashed_password",
		Role:     "user",
		PersonID: 2,
		Person: personModel.Model{
			Model: gorm.Model{ID: 2},
			Name:  "John Doe",
			Email: "john@example.com",
		},
	}

	domainUser := model.ToDomain()

	assert.NotNil(t, domainUser)
	assert.Equal(t, uint(1), domainUser.ID)
	assert.Equal(t, "user", domainUser.Role)
	assert.Equal(t, uint(2), domainUser.PersonID)
	assert.NotNil(t, domainUser.Password)
	assert.Equal(t, "hashed_password", domainUser.Password.GetHashed())
	assert.NotNil(t, domainUser.Person)
	assert.Equal(t, "John Doe", domainUser.Person.Name)
	assert.Equal(t, "john@example.com", domainUser.Person.Email)
	assert.Equal(t, now, domainUser.CreatedAt)
	assert.Nil(t, domainUser.DeletedAt)
}

func TestUserModel_FromDomain(t *testing.T) {
	now := time.Now()
	pass := domain.NewPasswordFromHash("hashed_password", nil)
	domainUser := &domain.User{
		ID:       1,
		Role:     "user",
		PersonID: 2,
		Person: &domain.Person{
			ID:    2,
			Name:  "John Doe",
			Email: "john@example.com",
		},
		Password:  &pass,
		CreatedAt: now,
		UpdatedAt: now,
	}

	model := &Model{}
	model.FromDomain(domainUser)

	assert.Equal(t, uint(1), model.ID)
	assert.Equal(t, "user", model.Role)
	assert.Equal(t, uint(2), model.PersonID)
	assert.Equal(t, "hashed_password", model.Password)
	assert.Equal(t, "John Doe", model.Person.Name)
	assert.Equal(t, "john@example.com", model.Person.Email)
	assert.Equal(t, now, model.CreatedAt)
	assert.False(t, model.DeletedAt.Valid)
}

func TestUserModel_FromDomain_Nil(t *testing.T) {
	model := &Model{Role: "admin"}
	model.FromDomain(nil)
	assert.Equal(t, "admin", model.Role)
}
