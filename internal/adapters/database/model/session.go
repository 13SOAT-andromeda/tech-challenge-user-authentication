package model

import "tech-challenge-user-validation/internal/core/ports"

type SessionModel struct {
	PK        string `dynamodbav:"pk"`
	UserID    string `dynamodbav:"user_id"`
	ExpiresAt int64  `dynamodbav:"expires_at"`
}

func (m *SessionModel) ToDomain() *ports.Session {
	return &ports.Session{
		ID:        m.PK,
		UserID:    m.UserID,
		ExpiresAt: m.ExpiresAt,
	}
}
