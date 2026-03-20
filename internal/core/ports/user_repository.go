package ports

import (
	"context"
	"tech-challenge-user-validation/internal/core/domain"
)

type UserSearch struct {
	Name string
}

type UserRepository interface {
	GetByDocument(ctx context.Context, document string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Search(ctx context.Context, params UserSearch) []domain.User
}
