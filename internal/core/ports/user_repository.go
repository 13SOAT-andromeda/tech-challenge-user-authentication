package ports

import (
	"context"
	"tech-challenge-user-validation/internal/core/domain"
)

type UserSearch struct {
	Name    string
	Email   string
	Contact string
}

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Search(ctx context.Context, params UserSearch) []domain.User
	GetByDocument(ctx context.Context, document string) (*domain.User, error) // Keeping this as it's heavily used in AuthUseCase, though not in the reference User repo
}
