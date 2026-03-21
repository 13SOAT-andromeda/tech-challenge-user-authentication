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

	accessExpiry := 1 * time.Hour

	userOutput := ports.UserOutput{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Contact: user.Contact,
		Role:    user.Role,
	}

	return &ports.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessExpiry.Seconds()),
		JTI:          jti,
		User:         userOutput,
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

func (uc *AuthUseCase) Validate(ctx context.Context, tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return uc.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return false, errors.New("invalid token")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("invalid token claims")
	}

	jti, ok := mapClaims["jti"].(string)
	if !ok || jti == "" {
		return false, errors.New("invalid token: missing jti")
	}

	session, err := uc.sessionService.GetByID(ctx, jti)
	if err != nil {
		return false, err
	}
	if session == nil {
		return false, errors.New("session not found or revoked")
	}

	return true, nil
}
