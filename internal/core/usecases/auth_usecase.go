package usecases

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"time"

	"tech-challenge-user-validation/internal/core/ports"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCase struct {
	userRepo       ports.UserRepository
	tokenRepo      ports.TokenRepository
	sessionService ports.SessionService
	jwtService     ports.JWTService
	jwtSecret      []byte
}

func NewAuthUseCase(
	userRepo ports.UserRepository,
	sessionService ports.SessionService,
	jwtService ports.JWTService,
	secret string,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		sessionService: sessionService,
		jwtService:     jwtService,
		jwtSecret:      []byte(secret),
	}
}

var DocumentRegex = regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)

func (uc *AuthUseCase) Login(ctx context.Context, input ports.LoginInput) (*ports.LoginOutput, error) {
	if !DocumentRegex.MatchString(input.Document) {
		return nil, errors.New("invalid document format")
	}

	user, err := uc.userRepo.GetByDocument(ctx, input.Document)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if err := user.Password.Compare(input.Password); err != nil || user.DeletedAt != nil || !user.IsActive {
		return nil, errors.New("user not found")
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshExpiry := 7 * 24 * time.Hour
	sessionExpiresAt := time.Now().Add(refreshExpiry)

	session, err := uc.sessionService.Create(ctx, , refreshToken, sessionExpiresAt)
	if err != nil {
		return nil, err
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role, session.ID)
	if err != nil {
		return nil, err
	}

	accessExpiry := 1 * time.Hour

	userOutput := ports.UserOutput{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Contact: user.Contact,
		Role:    user.Role,
	}

	output := &ports.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessExpiry.Seconds()),
		User:         userOutput,
	}

	return output, nil
}
