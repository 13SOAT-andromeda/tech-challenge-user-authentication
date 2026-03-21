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
	ValidateToken(tokenString string) (*JWTClaims, error)
	ExtractUserIDFromToken(tokenString string) (uint, error)
	IsTokenExpired(tokenString string) bool
	RefreshAccessToken(refreshTokenString, email, role, sessionID string) (string, error)
}
