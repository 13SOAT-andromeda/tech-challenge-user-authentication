package ports

import "time"

type JWTClaims struct {
	UserID    uint
	JTI       string
	Email     string
	Role      string
	SessionID string
	ExpiresAt time.Time
}

type JWTService interface {
	GenerateAccessToken(userID uint, email, role string, sessionID string) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
}
