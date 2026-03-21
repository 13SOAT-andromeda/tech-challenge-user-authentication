package model

import (
	"time"
)

type SessionModel struct {
	SessionID string    `dynamodbav:"token_id"`
	UserID    string    `dynamodbav:"user_id"`
	ExpiresAt time.Time `dynamodbav:"expires_at,omitempty"`
}
