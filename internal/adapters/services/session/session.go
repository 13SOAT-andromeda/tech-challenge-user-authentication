package session

import (
	"context"
	"errors"
	"strings"
	"time"

	"tech-challenge-user-validation/internal/core/ports"
)

type repository interface {
	Create(ctx context.Context, s *ports.Session) (*ports.Session, error)
	FindByID(ctx context.Context, sessionID uint) (*ports.Session, error)
}

type sessionService struct {
	repo repository
}

func NewSessionService(repo repository) *sessionService {
	return &sessionService{repo: repo}
}

func (s *sessionService) Create(ctx context.Context, sessionID string, userID string, expiresAt int64) (*ports.Session, error) {
	if sessionID == "" {
		return nil, errors.New("invalid session ID")
	}
	if userID == "" {
		return nil, errors.New("invalid user ID")
	}
	if expiresAt <= time.Now().Unix() {
		return nil, errors.New("expiration date must be in the future")
	}
	session := &ports.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return s.repo.Create(ctx, session)
}

func (s *sessionService) GetByID(ctx context.Context, sessionID string) (*ports.Session, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil, errors.New("invalid session ID")
	}
	// adapter stores uint IDs; this adapter's repo uses uint keys
	// fall through to repo
	return nil, errors.New("not implemented")
}
