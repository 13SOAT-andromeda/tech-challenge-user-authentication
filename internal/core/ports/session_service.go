package ports

import "context"

type Session struct {
	ID        string
	UserID    string
	ExpiresAt int64
}

type SessionService interface {
	Create(ctx context.Context, sessionID string, userID string, expiresAt int64) (*Session, error)
	GetByID(ctx context.Context, sessionID string) (*Session, error)
}
