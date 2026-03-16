package ports

import (
	"context"
	"tech-challenge-user-validation/internal/core/domain"
)

type UserRepository interface {
	GetByDocument(ctx context.Context, Document string) (*domain.User, error)
}
