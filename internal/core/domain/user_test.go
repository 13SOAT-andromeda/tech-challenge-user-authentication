package domain

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	now := time.Now()
	user := User{
		ID:        1,
		Document:  "123.456.789-00",
		IsActive:  true,
		CreatedAt: now,
	}

	if user.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}
	if user.Document != "123.456.789-00" {
		t.Errorf("expected Document 123.456.789-00, got %s", user.Document)
	}
	if !user.IsActive {
		t.Errorf("expected IsActive true, got %v", user.IsActive)
	}
	if !user.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, user.CreatedAt)
	}
}
