package usecases

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"tech-challenge-user-validation/internal/core/ports"
)

type AuthUseCase struct {
	userRepo       ports.UserRepository
	tokenRepo      ports.TokenRepository
	sessionService ports.SessionService
	jwtService     ports.JWTService
}

func NewAuthUseCase(
	userRepo ports.UserRepository,
	tokenRepo ports.TokenRepository,
	sessionService ports.SessionService,
	jwtService ports.JWTService,
	secret string,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		sessionService: sessionService,
		jwtService:     jwtService,
	}
}

var DocumentRegex = regexp.MustCompile(`^\d{11}$`)
var nonDigit = regexp.MustCompile(`\D`)

func normalizeDocument(doc string) string {
	return nonDigit.ReplaceAllString(doc, "")
}

func (uc *AuthUseCase) Login(ctx context.Context, input ports.LoginInput) (*ports.LoginOutput, error) {
	input.Document = normalizeDocument(input.Document)
	if !DocumentRegex.MatchString(input.Document) {
		return nil, errors.New("invalid document format")
	}

	user, err := uc.userRepo.GetByDocument(ctx, input.Document)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := user.Password.Compare(input.Password); err != nil || user.DeletedAt != nil || !user.IsActive {
		return nil, errors.New("invalid credentials")
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	claims, err := uc.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	jti := claims.JTI

	refreshExpiry := 7 * 24 * time.Hour
	sessionExpiresAt := time.Now().Add(refreshExpiry)

	session, err := uc.sessionService.Create(ctx, jti, strconv.Itoa(int(user.ID)), sessionExpiresAt.Unix())
	if err != nil {
		return nil, err
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role, session.ID)
	if err != nil {
		return nil, err
	}

	return &ports.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *AuthUseCase) Refresh(ctx context.Context, input ports.RefreshInput) (*ports.RefreshOutput, error) {
	claims, err := uc.jwtService.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	session, err := uc.sessionService.GetByID(ctx, claims.JTI)
	if err != nil || session == nil {
		return nil, errors.New("session not found or revoked")
	}

	if time.Now().Unix() > session.ExpiresAt {
		return nil, errors.New("session expired")
	}

	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role, session.ID)
	if err != nil {
		return nil, err
	}

	accessExpiry := 1 * time.Hour
	return &ports.RefreshOutput{
		AccessToken: accessToken,
		ExpiresIn:   int64(accessExpiry.Seconds()),
	}, nil
}

func (uc *AuthUseCase) Logout(ctx context.Context, tokenString string) error {
	claims, err := uc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return errors.New("invalid token")
	}

	sessionID := claims.SessionID
	if sessionID == "" {
		sessionID = claims.JTI
	}

	return uc.sessionService.Delete(ctx, sessionID)
}
