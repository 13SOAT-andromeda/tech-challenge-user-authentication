package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"tech-challenge-user-validation/internal/core/domain"
	"tech-challenge-user-validation/internal/core/ports"

)

type mockHasher struct{}

func (h *mockHasher) Hash(password string) (string, error) { return password, nil }
func (h *mockHasher) Compare(hashedPassword, password string) error {
	if hashedPassword != password {
		return errors.New("invalid password")
	}
	return nil
}

type mockUserRepository struct {
	getUserFunc    func(ctx context.Context, document string) (*domain.User, error)
	getUserByIDFunc func(ctx context.Context, id uint) (*domain.User, error)
}

func (m *mockUserRepository) GetByDocument(ctx context.Context, document string) (*domain.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, document)
	}
	return nil, nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, nil
}

type mockTokenRepository struct {
	saveFunc func(ctx context.Context, pk string, token string, expiresAt int64) error
}

func (m *mockTokenRepository) Save(ctx context.Context, pk string, token string, expiresAt int64) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, pk, token, expiresAt)
	}
	return nil
}

type mockSessionService struct {
	createFunc func(ctx context.Context, sessionID string, userID string, expiresAt int64) (*ports.Session, error)
	getFunc    func(ctx context.Context, sessionID string) (*ports.Session, error)
	deleteFunc func(ctx context.Context, sessionID string) error
}

func (m *mockSessionService) Create(ctx context.Context, sessionID string, userID string, expiresAt int64) (*ports.Session, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, sessionID, userID, expiresAt)
	}
	return &ports.Session{ID: sessionID, UserID: userID, ExpiresAt: expiresAt}, nil
}

func (m *mockSessionService) GetByID(ctx context.Context, sessionID string) (*ports.Session, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, sessionID)
	}
	return nil, nil
}

func (m *mockSessionService) Delete(ctx context.Context, sessionID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, sessionID)
	}
	return nil
}

type mockJWTService struct {
	generateAccessTokenFunc    func(userID uint, email, role, sessionID string) (string, error)
	generateRefreshTokenFunc   func(userID uint) (string, error)
	validateTokenFunc          func(tokenString string) (*ports.JWTClaims, error)
	validateRefreshTokenFunc   func(tokenString string) (*ports.JWTClaims, error)
	extractUserIDFunc          func(tokenString string) (uint, error)
	isTokenExpiredFunc         func(tokenString string) bool
	refreshAccessTokenFunc     func(refreshTokenString, email, role, sessionID string) (string, error)
}

func (m *mockJWTService) GenerateAccessToken(userID uint, email, role, sessionID string) (string, error) {
	if m.generateAccessTokenFunc != nil {
		return m.generateAccessTokenFunc(userID, email, role, sessionID)
	}
	return "access-token", nil
}

func (m *mockJWTService) GenerateRefreshToken(userID uint) (string, error) {
	if m.generateRefreshTokenFunc != nil {
		return m.generateRefreshTokenFunc(userID)
	}
	return "refresh-token", nil
}

func (m *mockJWTService) ValidateToken(tokenString string) (*ports.JWTClaims, error) {
	if m.validateTokenFunc != nil {
		return m.validateTokenFunc(tokenString)
	}
	return &ports.JWTClaims{
		UserID:    1,
		JTI:       "jti-default",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

func (m *mockJWTService) ValidateRefreshToken(tokenString string) (*ports.JWTClaims, error) {
	if m.validateRefreshTokenFunc != nil {
		return m.validateRefreshTokenFunc(tokenString)
	}
	return &ports.JWTClaims{
		UserID:    1,
		JTI:       "jti-default",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

func (m *mockJWTService) ExtractUserIDFromToken(tokenString string) (uint, error) {
	if m.extractUserIDFunc != nil {
		return m.extractUserIDFunc(tokenString)
	}
	return 1, nil
}

func (m *mockJWTService) IsTokenExpired(tokenString string) bool {
	if m.isTokenExpiredFunc != nil {
		return m.isTokenExpiredFunc(tokenString)
	}
	return false
}

func (m *mockJWTService) RefreshAccessToken(refreshTokenString, email, role, sessionID string) (string, error) {
	if m.refreshAccessTokenFunc != nil {
		return m.refreshAccessTokenFunc(refreshTokenString, email, role, sessionID)
	}
	return "access-token", nil
}

func makeUser(document, rawPassword string, active bool) *domain.User {
	pass := domain.NewPasswordFromHash(rawPassword, &mockHasher{})
	return &domain.User{
		ID:       1,
		Role:     "user",
		PersonID: 1,
		Person: &domain.Person{
			Name:     "Barbara",
			Email:    "barbara@exemplo.com",
			Contact:  "11999999999",
			Document: document,
			IsActive: active,
		},
		Password: &pass,
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	ctx := context.Background()
	secret := "secret"

	t.Run("should fail with invalid document format", func(t *testing.T) {
		uc := NewAuthUseCase(nil, nil, nil, &mockJWTService{}, secret)
		_, err := uc.Login(ctx, ports.LoginInput{
			Document: "invalid",
			Password: "123456",
		})
		if err == nil || err.Error() != "invalid document format" {
			t.Fatalf("expected 'invalid document format', got: %v", err)
		}
	})

	t.Run("should fail if user not found", func(t *testing.T) {
		userRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, document string) (*domain.User, error) {
				return nil, nil
			},
		}
		uc := NewAuthUseCase(userRepo, nil, nil, &mockJWTService{}, secret)
		_, err := uc.Login(ctx, ports.LoginInput{
			Document: "123.456.789-00",
			Password: "123456",
		})
		if err == nil || err.Error() != "user not found" {
			t.Fatalf("expected 'user not found', got: %v", err)
		}
	})

	t.Run("should fail if user is inactive", func(t *testing.T) {
		userRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, document string) (*domain.User, error) {
				return makeUser(document, "123456", false), nil
			},
		}
		uc := NewAuthUseCase(userRepo, nil, nil, &mockJWTService{}, secret)
		_, err := uc.Login(ctx, ports.LoginInput{
			Document: "123.456.789-00",
			Password: "123456",
		})
		if err == nil || err.Error() != "invalid credentials" {
			t.Fatalf("expected 'invalid credentials' for inactive user, got: %v", err)
		}
	})

	t.Run("should fail if password is invalid", func(t *testing.T) {
		userRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, document string) (*domain.User, error) {
				return makeUser(document, "correct-password", true), nil
			},
		}
		uc := NewAuthUseCase(userRepo, nil, nil, &mockJWTService{}, secret)
		_, err := uc.Login(ctx, ports.LoginInput{
			Document: "123.456.789-00",
			Password: "wrong-password",
		})
		if err == nil || err.Error() != "invalid credentials" {
			t.Fatalf("expected 'invalid credentials' for invalid password, got: %v", err)
		}
	})

	t.Run("should succeed and return access_token/refresh_token/jti", func(t *testing.T) {
		expectedJTI := "session-jti-123"
		exp := time.Now().Add(24 * time.Hour)

		userRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, document string) (*domain.User, error) {
				return makeUser(document, "123456", true), nil
			},
		}

		var capturedSessionID string
		var capturedUserID string
		sessionSvc := &mockSessionService{
			createFunc: func(ctx context.Context, sessionID string, userID string, expiresAt int64) (*ports.Session, error) {
				capturedSessionID = sessionID
				capturedUserID = userID
				return &ports.Session{ID: sessionID, UserID: userID, ExpiresAt: expiresAt}, nil
			},
		}

		jwtSvc := &mockJWTService{
			generateRefreshTokenFunc: func(userID uint) (string, error) {
				return "refresh-token-abc", nil
			},
			validateRefreshTokenFunc: func(tokenString string) (*ports.JWTClaims, error) {
				return &ports.JWTClaims{
					UserID:    1,
					JTI:       expectedJTI,
					Email:     "barbara@exemplo.com",
					Role:      "user",
					SessionID: expectedJTI,
					ExpiresAt: exp,
				}, nil
			},
			generateAccessTokenFunc: func(userID uint, email, role, sessionID string) (string, error) {
				if sessionID != expectedJTI {
					t.Fatalf("expected sessionID %s, got %s", expectedJTI, sessionID)
				}
				return "access-token-xyz", nil
			},
		}

		uc := NewAuthUseCase(userRepo, &mockTokenRepository{}, sessionSvc, jwtSvc, secret)

		out, err := uc.Login(ctx, ports.LoginInput{
			Document: "123.456.789-00",
			Password: "123456",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if out.AccessToken == "" {
			t.Fatal("expected access token")
		}
		if out.RefreshToken == "" {
			t.Fatal("expected refresh token")
		}
		if capturedSessionID != expectedJTI {
			t.Fatalf("expected session create with jti %s, got %s", expectedJTI, capturedSessionID)
		}
		if capturedUserID != "1" {
			t.Fatalf("expected session create with userID '1', got %s", capturedUserID)
		}
	})
}

