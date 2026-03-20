package session

import (
	"context"
	"errors"
	"strconv"
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

func (s *sessionService) Create(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) (*ports.Session, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, errors.New("refresh token cannot be empty")
	}

	now := time.Now()
	if !expiresAt.After(now) {
		return nil, errors.New("expiration date must be in the future")
	}

	session := &ports.Session{
		UserID:    strconv.Itoa(int(userID)),
		ExpiresAt: expiresAt,
	}
	return s.repo.Create(ctx, session)
}

func (s *sessionService) GetByID(ctx context.Context, sessionID string) (*ports.Session, error) {
	id, err := strconv.ParseUint(sessionID, 10, 64)
	if err != nil || id == 0 {
		return nil, errors.New("invalid session ID")
	}

	return s.repo.FindByID(ctx, uint(id))
}
