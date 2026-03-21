package ports

import (
	"context"
	"tech-challenge-user-validation/internal/core/domain"
)

type UserRepository interface {
	GetByDocument(ctx context.Context, document string) (*domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
}
