package usecases

import (
	"context"
	"errors"
	"regexp"
	"time"

	"tech-challenge-user-validation/internal/core/ports"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	userRepo    ports.UserRepository
	tokenRepo   ports.TokenRepository
	sessionService ports.SessionService
	jwtSecret   []byte
}

func NewAuthUseCase(userRepo ports.UserRepository, tokenRepo ports.TokenRepository, sessionService ports.SessionService, secret string) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		sessionService: sessionService,
		jwtSecret:      []byte(secret),
	}
}

var DocumentRegex = regexp.MustCompile(`^(\d{3}\.\d{3}\.\d{3}\-\d{2})?$`)

func (uc *AuthUseCase) Authenticate(ctx context.Context, Document string) (string, error) {
	if !DocumentRegex.MatchString(Document) {
		return "", errors.New("invalid Document format")
	}

	user, err := uc.userRepo.GetByDocument(ctx, Document)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	if !user.IsActive {
		return "", errors.New("user is inactive")
	}

	// Generate JTI (UUID v4)
	jti := uuid.New().String()

	// Expiry: 24 hours
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"Document": user.Document,
		"jti":      jti,
		"exp":      expiresAt,
	})

	tokenString, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return "", err
	}

	// 1. Save minimalist session in SessionService (DynamoDB)
	_, err = uc.sessionService.Create(ctx, jti, user.Document, expiresAt)
	if err != nil {
		return "", err
	}

	// 2. Backward compatibility: Save token in TokenRepository (if still needed)
	err = uc.tokenRepo.Save(ctx, user.Document, tokenString, expiresAt)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (uc *AuthUseCase) Validate(ctx context.Context, tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return uc.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return false, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("invalid claims")
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return false, errors.New("jti not found in token")
	}

	// Verify session exists in DynamoDB
	session, err := uc.sessionService.GetByID(ctx, jti)
	if err != nil {
		return false, err
	}

	if session == nil {
		return false, errors.New("session not found or revoked")
	}

	return true, nil
}
