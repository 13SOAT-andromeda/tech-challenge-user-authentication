package services

import (
	"context"
	"errors"
	"time"

	"tech-challenge-user-validation/internal/core/ports"
)

type sessionService struct{}

func NewSessionService() *sessionService {
	return &sessionService{}
}

func (s *sessionService) Create(ctx context.Context, sessionID string, userID string, expiresAt int64) (*ports.Session, error) {
	if sessionID == "" {
		return nil, errors.New("invalid session ID")
	}
	if userID == "" {
		return nil, errors.New("invalid user ID")
	}
	if expiresAt <= time.Now().Unix() {
		return nil, errors.New("expiration date cannot be in the past")
	}
	return &ports.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *sessionService) GetByID(ctx context.Context, sessionID string) (*ports.Session, error) {
	return nil, errors.New("not implemented")
}
