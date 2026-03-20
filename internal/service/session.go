package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"tech-challenge-user-validation/internal/core/ports"

	"github.com/google/uuid"
)

type sessionService struct{}

func NewSessionService() *sessionService {
	return &sessionService{}
}

func (s *sessionService) Create(ctx context.Context, userID uint, expiresAt time.Time) (*ports.Session, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	if expiresAt.Before(time.Now()) {
		return nil, errors.New("expiration date cannot be in the past")
	}

	return &ports.Session{
		ID:        uuid.NewString(),
		UserID:    strconv.Itoa(int(userID)),
		ExpiresAt: time.time,
	}, nil
}

func (s *sessionService) GetByID(ctx context.Context, sessionID string) (*ports.Session, error) {
	return nil, errors.New("not implemented")
}
