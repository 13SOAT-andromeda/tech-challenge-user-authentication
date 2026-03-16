package usecases

import (
	"context"
	"errors"
	"testing"
	"tech-challenge-user-validation/internal/core/domain"
)

type mockUserRepository struct {
	getUserFunc func(ctx context.Context, cpf string) (*domain.User, error)
}

func (m *mockUserRepository) GetByCPF(ctx context.Context, cpf string) (*domain.User, error) {
	return m.getUserFunc(ctx, cpf)
}

type mockTokenRepository struct {
	saveFunc func(ctx context.Context, token *domain.Token) error
}

func (m *mockTokenRepository) Save(ctx context.Context, token *domain.Token) error {
	return m.saveFunc(ctx, token)
}

func TestAuthUseCase_Authenticate(t *testing.T) {
	ctx := context.Background()

	t.Run("should fail with invalid CPF format", func(t *testing.T) {
		uc := NewAuthUseCase(nil, nil, "secret")
		_, err := uc.Authenticate(ctx, "invalid-cpf")
		if err == nil {
			t.Fatal("expected error for invalid CPF format")
		}
	})

	t.Run("should fail if user not found", func(t *testing.T) {
		mockUserRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, cpf string) (*domain.User, error) {
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
			getUserFunc: func(ctx context.Context, cpf string) (*domain.User, error) {
				return &domain.User{ID: 1, CPF: cpf, IsActive: false}, nil
			},
		}
		uc := NewAuthUseCase(mockUserRepo, nil, "secret")
		_, err := uc.Authenticate(ctx, "123.456.789-00")
		if err == nil {
			t.Fatal("expected error when user is inactive")
		}
	})

	t.Run("should succeed and save token if user is active", func(t *testing.T) {
		mockUserRepo := &mockUserRepository{
			getUserFunc: func(ctx context.Context, cpf string) (*domain.User, error) {
				return &domain.User{ID: 1, CPF: cpf, IsActive: true}, nil
			},
		}
		tokenSaved := false
		mockTokenRepo := &mockTokenRepository{
			saveFunc: func(ctx context.Context, token *domain.Token) error {
				if token.UserID == 1 && token.TokenID != "" {
					tokenSaved = true
				}
				return nil
			},
		}
		uc := NewAuthUseCase(mockUserRepo, mockTokenRepo, "secret")
		token, err := uc.Authenticate(ctx, "123.456.789-00")
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
