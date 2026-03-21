package ports

import (
	"context"
	"tech-challenge-user-validation/internal/core/domain"
)

type UserRepository interface {
	GetByDocument(ctx context.Context, document string) (*domain.User, error)
}
