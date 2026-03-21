package jwt

import (
	"fmt"
	"time"

	"tech-challenge-user-validation/internal/core/ports"

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
	refreshSecret      []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewService(secret, refreshSecret string, accessTokenExpiry, refreshTokenExpiry time.Duration) *Service {
	return &Service{
		secret:             []byte(secret),
		refreshSecret:      []byte(refreshSecret),
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func (s *Service) GenerateAccessToken(userID uint, email, role string, sessionID string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		JTI:       sessionID,
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
	return token.SignedString(s.refreshSecret)
}

func (s *Service) ValidateToken(tokenString string) (*ports.JWTClaims, error) {
	return s.parseToken(tokenString, s.secret)
}

func (s *Service) ValidateRefreshToken(tokenString string) (*ports.JWTClaims, error) {
	return s.parseToken(tokenString, s.refreshSecret)
}

func (s *Service) parseToken(tokenString string, secret []byte) (*ports.JWTClaims, error) {
	token, err := jwtv5.ParseWithClaims(tokenString, &Claims{}, func(token *jwtv5.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return &ports.JWTClaims{
		UserID:    claims.UserID,
		JTI:       claims.JTI,
		Email:     claims.Email,
		Role:      claims.Role,
		SessionID: claims.SessionID,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

func (s *Service) ExtractUserIDFromToken(tokenString string) (uint, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func (s *Service) IsTokenExpired(tokenString string) bool {
	_, err := s.ValidateToken(tokenString)
	return err != nil
}

func (s *Service) RefreshAccessToken(refreshTokenString, email, role, sessionID string) (string, error) {
	claims, err := s.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", err
	}
	return s.GenerateAccessToken(claims.UserID, email, role, sessionID)
}
