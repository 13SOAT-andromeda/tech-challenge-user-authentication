package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"tech-challenge-user-validation/internal/core/domain"
)

func TestUserModel_ToDomain(t *testing.T) {
	now := time.Now()
	userModel := UserModel{
		ID:        1,
		Document:  "12345678900",
		IsActive:  true,
		CreatedAt: now,
	}

	expected := domain.User{
		ID:        1,
		Document:  "12345678900",
		IsActive:  true,
		CreatedAt: now,
	}

	result := userModel.ToDomain()

	assert.Equal(t, expected, result)
}
