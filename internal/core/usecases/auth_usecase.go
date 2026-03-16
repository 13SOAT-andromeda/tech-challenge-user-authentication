package usecases

import (
	"context"
	"errors"
	"regexp"
	"time"

	"tech-challenge-user-validation/internal/core/domain"
	"tech-challenge-user-validation/internal/core/ports"

	"github.com/golang-jwt/jwt/v5"
)

// Dummy import to satisfy golang-jwt/jwt/v5 if needed, but actually I need to check the exact usage.
// I'll use standard jwt-go style.

type AuthUseCase struct {
	userRepo  ports.UserRepository
	tokenRepo ports.TokenRepository
	jwtSecret []byte
}

func NewAuthUseCase(userRepo ports.UserRepository, tokenRepo ports.TokenRepository, secret string) *AuthUseCase {
	return &AuthUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtSecret: []byte(secret),
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

	if !user.IsActive {
		return "", errors.New("user is inactive")
	}

	// Generate JTI
	jti := time.Now().Format("20060102150405") // Simple JTI for this example

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"Document": user.Document,
		"jti": jti,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return "", err
	}

	// Save token in DynamoDB (integration logic)
	tokenEntity := &domain.Token{
		TokenID: jti,
		UserID:  user.ID,
	}

	err = uc.tokenRepo.Save(ctx, tokenEntity)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
