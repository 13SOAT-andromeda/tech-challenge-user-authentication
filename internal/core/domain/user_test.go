package domain

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	now := time.Now()
	user := User{
		ID:        1,
		PersonID:  2,
		CreatedAt: now,
	}

	if user.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}
	if user.PersonID != 2 {
		t.Errorf("expected PersonID 2, got %d", user.PersonID)
	}
	if !user.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, user.CreatedAt)
	}
}
