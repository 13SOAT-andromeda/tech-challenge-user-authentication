package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"tech-challenge-user-validation/internal/core/domain"
	"tech-challenge-user-validation/internal/core/usecases"

	"github.com/aws/aws-lambda-go/events"
)

type mockUserRepository struct{}

func (m *mockUserRepository) GetByDocument(ctx context.Context, Document string) (*domain.User, error) {
	if Document == "123.456.789-00" {
		return &domain.User{ID: 1, Document: Document, IsActive: true}, nil
	}
	return nil, nil
}

type mockTokenRepository struct{}

func (m *mockTokenRepository) Save(ctx context.Context, pk string, token string, expiresAt int64) error {
	return nil
}

func TestAuthHandler_Handle(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepository{}
	tokenRepo := &mockTokenRepository{}
	uc := usecases.NewAuthUseCase(userRepo, tokenRepo, "secret")
	h := NewAuthHandler(uc)

	t.Run("should return 400 if x-user-cpf header is missing", func(t *testing.T) {
		req := events.APIGatewayProxyRequest{}
		resp, err := h.Handle(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("should return 200 and token if Document is valid and user is active", func(t *testing.T) {
		req := events.APIGatewayProxyRequest{
			Headers: map[string]string{
				"x-user-cpf": "123.456.789-00",
			},
		}
		resp, err := h.Handle(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		var body map[string]string
		json.Unmarshal([]byte(resp.Body), &body)
		if body["token"] == "" {
			t.Error("expected token in response body")
		}
	})
}
