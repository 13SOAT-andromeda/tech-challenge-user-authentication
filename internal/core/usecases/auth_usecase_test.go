package usecases

import (
	"context"
	"testing"
	"time"
	"tech-challenge-user-validation/internal/core/domain"
	"tech-challenge-user-validation/internal/core/ports"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type mockUserRepository struct {
	getUserFunc func(ctx context.Context, Document string) (*domain.User, error)
}

func (m *mockUserRepository) GetByDocument(ctx context.Context, Document string) (*domain.User, error) {
	return m.getUserFunc(ctx, Document)
}

type mockTokenRepository struct {
	saveFunc func(ctx context.Context, pk string, token string, expiresAt int64) error
}

func (m *mockTokenRepository) Save(ctx context.Context, pk string, token string, expiresAt int64) error {
	return m.saveFunc(ctx, pk, token, expiresAt)
}

type mockSessionService struct {
	createFunc func(ctx context.Context, sessionID string, userID string, expiresAt int64) (*ports.Session, error)
	getFunc    func(ctx context.Context, sessionID string) (*ports.Session, error)
}

func (m *mockSessionService) Create(ctx context.Context, sessionID string, userID string, expiresAt int64) (*ports.Session, error) {
	return m.createFunc(ctx, sessionID, userID, expiresAt)
}

func (m *mockSessionService) GetByID(ctx context.Context, sessionID string) (*ports.Session, error) {
	return m.getFunc(ctx, sessionID)
}

func TestAuthUseCase_Authenticate(t *testing.T) {
	ctx := context.Background()
	secret := "secret"

	t.Run("should fail with invalid Document format", func(t *testing.T) {
		uc := NewAuthUseCase(nil, nil, nil, secret)
		_, err := uc.Authenticate(ctx, "invalid-Document")
		if err == nil {
			t.Fatal("expected error for invalid Document format")
		}
	})

	t.Run("should fail if user not found", func(t *testing.T) {
		mockUserRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, Document string) (*domain.User, error) {
				return nil, nil // Return nil for not found
			},
		}
		uc := NewAuthUseCase(mockUserRepo, nil, nil, secret)
		_, err := uc.Authenticate(ctx, "123.456.789-00")
		if err == nil || err.Error() != "user not found" {
			t.Fatalf("expected 'user not found' error, got: %v", err)
		}
	})

	t.Run("should fail if user is inactive", func(t *testing.T) {
		mockUserRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, Document string) (*domain.User, error) {
				return &domain.User{ID: 1, Document: Document, IsActive: false}, nil
			},
		}
		uc := NewAuthUseCase(mockUserRepo, nil, nil, secret)
		_, err := uc.Authenticate(ctx, "123.456.789-00")
		if err == nil || err.Error() != "user is inactive" {
			t.Fatal("expected error when user is inactive")
		}
	})

	t.Run("should succeed and save token/session if user is active", func(t *testing.T) {
		document := "123.456.789-00"
		mockUserRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, Document string) (*domain.User, error) {
				return &domain.User{ID: 1, Document: document, IsActive: true}, nil
			},
		}
		
		tokenSaved := false
		mockTokenRepo := &mockTokenRepository{
			saveFunc: func(ctx context.Context, pk string, token string, expiresAt int64) error {
				tokenSaved = true
				return nil
			},
		}

		sessionCreated := false
		var capturedJTI string
		mockSessionSvc := &mockSessionService{
			createFunc: func(ctx context.Context, sessionID string, userID string, expiresAt int64) (*ports.Session, error) {
				sessionCreated = true
				capturedJTI = sessionID
				return &ports.Session{ID: sessionID}, nil
			},
		}

		uc := NewAuthUseCase(mockUserRepo, mockTokenRepo, mockSessionSvc, secret)
		tokenString, err := uc.Authenticate(ctx, document)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if tokenString == "" {
			t.Fatal("expected token string")
		}

		// Verify JTI in token
		token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		claims := token.Claims.(jwt.MapClaims)
		jti := claims["jti"].(string)

		if _, err := uuid.Parse(jti); err != nil {
			t.Fatalf("jti is not a valid uuid: %v", err)
		}

		if jti != capturedJTI {
			t.Fatal("jti in token does not match jti in session service")
		}

		if !tokenSaved {
			t.Fatal("expected token to be saved in token repository")
		}

		if !sessionCreated {
			t.Fatal("expected session to be created in session service")
		}
	})
}

func TestAuthUseCase_Validate(t *testing.T) {
	ctx := context.Background()
	secret := "secret"
	jti := uuid.New().String()
	
	// Generate a valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": jti,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	t.Run("should succeed if token is valid and session exists", func(t *testing.T) {
		mockSessionSvc := &mockSessionService{
			getFunc: func(ctx context.Context, sessionID string) (*ports.Session, error) {
				if sessionID == jti {
					return &ports.Session{ID: jti}, nil
				}
				return nil, nil
			},
		}
		uc := NewAuthUseCase(nil, nil, mockSessionSvc, secret)
		valid, err := uc.Validate(ctx, tokenString)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !valid {
			t.Fatal("expected token to be valid")
		}
	})

	t.Run("should fail if session does not exist", func(t *testing.T) {
		mockSessionSvc := &mockSessionService{
			getFunc: func(ctx context.Context, sessionID string) (*ports.Session, error) {
				return nil, nil
			},
		}
		uc := NewAuthUseCase(nil, nil, mockSessionSvc, secret)
		valid, err := uc.Validate(ctx, tokenString)
		if err == nil || err.Error() != "session not found or revoked" {
			t.Fatalf("expected 'session not found or revoked' error, got: %v", err)
		}
		if valid {
			t.Fatal("expected token to be invalid")
		}
	})
}
