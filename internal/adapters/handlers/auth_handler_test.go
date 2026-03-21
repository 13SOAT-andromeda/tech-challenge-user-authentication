package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"tech-challenge-user-validation/internal/core/domain"
	"tech-challenge-user-validation/internal/core/ports"
	"tech-challenge-user-validation/internal/core/usecases"

	"github.com/aws/aws-lambda-go/events"
)

type plainHasher struct{}

func (h *plainHasher) Hash(password string) (string, error) { return password, nil }
func (h *plainHasher) Compare(hashedPassword, password string) error {
	if hashedPassword != password {
		return errors.New("invalid password")
	}
	return nil
}

type mockUserRepository struct{}

func (m *mockUserRepository) GetByDocument(ctx context.Context, document string) (*domain.User, error) {
	if document == "123.456.789-00" {
		return &domain.User{
			ID:       1,
			Name:     "Barbara",
			Email:    "barbara@exemplo.com",
			Contact:  "11999999999",
			Role:     "user",
			Document: document,
			IsActive: true,
			Password: domain.NewPasswordFromHash("123456", &plainHasher{}),
		}, nil
	}
	return nil, nil
}


type mockTokenRepository struct{}

func (m *mockTokenRepository) Save(ctx context.Context, pk string, token string, expiresAt int64) error {
	return nil
}

type mockSessionService struct{}

func (m *mockSessionService) Create(ctx context.Context, sessionID string, userID string, expiresAt int64) (*ports.Session, error) {
	return &ports.Session{ID: sessionID, UserID: userID, ExpiresAt: expiresAt}, nil
}

func (m *mockSessionService) GetByID(ctx context.Context, sessionID string) (*ports.Session, error) {
	return nil, nil
}

type mockJWTService struct{}

func (m *mockJWTService) GenerateAccessToken(userID uint, email, role, sessionID string) (string, error) {
	return "access-token-123", nil
}

func (m *mockJWTService) GenerateRefreshToken(userID uint) (string, error) {
	return "refresh-token-123", nil
}

func (m *mockJWTService) ValidateToken(tokenString string) (*ports.JWTClaims, error) {
	return &ports.JWTClaims{
		UserID:    1,
		JTI:       "jti-123",
		Email:     "barbara@exemplo.com",
		Role:      "user",
		SessionID: "jti-123",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

func (m *mockJWTService) ExtractUserIDFromToken(tokenString string) (uint, error) {
	return 1, nil
}

func (m *mockJWTService) IsTokenExpired(tokenString string) bool {
	return false
}

func (m *mockJWTService) RefreshAccessToken(refreshTokenString, email, role, sessionID string) (string, error) {
	return "access-token-123", nil
}

func TestAuthHandler_Handle(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepository{}
	tokenRepo := &mockTokenRepository{}
	sessionSvc := &mockSessionService{}
	jwtSvc := &mockJWTService{}

	uc := usecases.NewAuthUseCase(userRepo, tokenRepo, sessionSvc, jwtSvc, "secret")
	h := NewAuthHandler(uc)

	t.Run("should return 400 if request body is invalid", func(t *testing.T) {
		req := events.APIGatewayProxyRequest{Body: ""}
		resp, err := h.Handle(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("should return 200 and access_token/refresh_token/jti", func(t *testing.T) {
		req := events.APIGatewayProxyRequest{
			Body: `{"document":"123.456.789-00","password":"123456"}`,
		}

		resp, err := h.Handle(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var body ports.LoginOutput
		if err := json.Unmarshal([]byte(resp.Body), &body); err != nil {
			t.Fatalf("failed to parse body: %v", err)
		}

		if body.AccessToken == "" {
			t.Error("expected access_token in response body")
		}
		if body.RefreshToken == "" {
			t.Error("expected refresh_token in response body")
		}
		if body.JTI == "" {
			t.Error("expected jti in response body")
		}
	})
}
