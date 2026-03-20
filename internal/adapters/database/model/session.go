package model

import (
	"time"
)

type SessionModel struct {
	UserID       uint      `dynamodbav:"token_id"`
	RefreshToken *string   `dynamodbav:"refresh_token,omitempty"`
	ExpiresAt    time.Time `dynamodbav:"expires_at,omitempty"`
}
