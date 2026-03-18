package usecases

import (
	"context"
	"errors"
	"testing"
	"tech-challenge-user-validation/internal/core/domain"
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

func TestAuthUseCase_Authenticate(t *testing.T) {
	ctx := context.Background()

	t.Run("should fail with invalid Document format", func(t *testing.T) {
		uc := NewAuthUseCase(nil, nil, "secret")
		_, err := uc.Authenticate(ctx, "invalid-Document")
		if err == nil {
			t.Fatal("expected error for invalid Document format")
		}
	})

	t.Run("should fail if user not found", func(t *testing.T) {
		mockUserRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, Document string) (*domain.User, error) {
				return nil, errors.New("user not found")
			},
		}
		uc := NewAuthUseCase(mockUserRepo, nil, "secret")
		_, err := uc.Authenticate(ctx, "123.456.789-00")
		if err == nil {
			t.Fatal("expected error when user not found")
		}
	})

	t.Run("should fail if user is inactive", func(t *testing.T) {
		mockUserRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, Document string) (*domain.User, error) {
				return &domain.User{ID: 1, Document: Document, IsActive: false}, nil
			},
		}
		uc := NewAuthUseCase(mockUserRepo, nil, "secret")
		_, err := uc.Authenticate(ctx, "123.456.789-00")
		if err == nil {
			t.Fatal("expected error when user is inactive")
		}
	})

	t.Run("should succeed and save token if user is active", func(t *testing.T) {
		document := "123.456.789-00"
		mockUserRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, Document string) (*domain.User, error) {
				return &domain.User{ID: 1, Document: document, IsActive: true}, nil
			},
		}
		tokenSaved := false
		mockTokenRepo := &mockTokenRepository{
			saveFunc: func(ctx context.Context, pk string, token string, expiresAt int64) error {
				if pk == document && token != "" && expiresAt > 0 {
					tokenSaved = true
				}
				return nil
			},
		}
		uc := NewAuthUseCase(mockUserRepo, mockTokenRepo, "secret")
		token, err := uc.Authenticate(ctx, document)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token == "" {
			t.Fatal("expected token string")
		}
		if !tokenSaved {
			t.Fatal("expected token to be saved")
		}
	})
}
