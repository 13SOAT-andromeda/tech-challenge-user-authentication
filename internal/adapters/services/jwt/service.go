package jwt

import (
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID    uint   `json:"user_id"`
	JTI       string `json:"jti"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	SessionID string `json:"session_id"`
	jwtv5.RegisteredClaims
}

type Service struct {
	secret             []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewService(secret string, accessTokenExpiry, refreshTokenExpiry time.Duration) *Service {
	return &Service{
		secret:             []byte(secret),
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func (s *Service) GenerateAccessToken(userID uint, email, role string, sessionID string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		JTI:       sessionID, // reutiliza jti da sessao
		Email:     email,
		Role:      role,
		SessionID: sessionID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(s.accessTokenExpiry)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
			NotBefore: jwtv5.NewNumericDate(time.Now()),
			Issuer:    "tech-challenge-s1",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) GenerateRefreshToken(userID uint) (string, error) {
	jti := uuid.NewString()
	claims := &Claims{
		UserID: userID,
		JTI:    jti,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(s.refreshTokenExpiry)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
			NotBefore: jwtv5.NewNumericDate(time.Now()),
			Issuer:    "tech-challenge-s1",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
