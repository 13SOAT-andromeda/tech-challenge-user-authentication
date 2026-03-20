package session

import (
	"context"
	"errors"
	"strings"
	"time"

	"tech-challenge-user-validation/internal/core/domain"
)

type sessionService struct {
	repo repository
}

func NewSessionService(repo repository) *sessionService {
	return &sessionService{repo: repo}
}

func (s *sessionService) Create(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) (*domain.Session, error) {
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

	session := (userID, refreshToken, expiresAt)
	return s.repo.Create(ctx, session)
}

func (s *sessionService) GetByID(ctx context.Context, sessionID uint) (*domain.Session, error) {
	if sessionID == 0 {
		return nil, errors.New("invalid session ID")
	}

	return s.repo.FindByID(ctx, sessionID)
}
