package repositories

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type DummyModel struct {
	ID int
}

func TestNewBaseRepository(t *testing.T) {
	db := &gorm.DB{}
	repo := NewBaseRepository[DummyModel](db)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}
