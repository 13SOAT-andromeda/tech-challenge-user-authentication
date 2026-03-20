package ports

import (
	"context"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type SessionService interface {
	Create(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) (*Session, error)
	GetByID(ctx context.Context, sessionID string) (*Session, error)
}
