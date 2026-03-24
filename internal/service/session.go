package services

import (
	"context"
	"errors"
	"time"

	"tech-challenge-user-validation/internal/adapters/database/model"
	"tech-challenge-user-validation/internal/core/ports"
)

type sessionRepository interface {
	Save(ctx context.Context, s model.SessionModel) error
	FindBySessionID(ctx context.Context, sessionID string) (*model.SessionModel, error)
	DeleteBySessionID(ctx context.Context, sessionID string) error
}

type sessionService struct {
	repo sessionRepository
}

func NewSessionService(repo sessionRepository) *sessionService {
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
		return nil, errors.New("expiration date cannot be in the past")
	}

	session := model.SessionModel{
		SessionID: sessionID,
		UserID:    userID,
		ExpiresAt: time.Unix(expiresAt, 0).UTC(),
	}
	if err := s.repo.Save(ctx, session); err != nil {
		return nil, err
	}

	return &ports.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *sessionService) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return errors.New("invalid session ID")
	}
	return s.repo.DeleteBySessionID(ctx, sessionID)
}

func (s *sessionService) GetByID(ctx context.Context, sessionID string) (*ports.Session, error) {
	if sessionID == "" {
		return nil, errors.New("invalid session ID")
	}

	m, err := s.repo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	return &ports.Session{
		ID:        m.SessionID,
		UserID:    m.UserID,
		ExpiresAt: m.ExpiresAt.Unix(),
	}, nil
}
