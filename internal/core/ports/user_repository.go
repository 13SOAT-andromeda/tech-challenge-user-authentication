package ports

import (
	"context"
	"tech-challenge-user-validation/internal/core/domain"
)

type UserRepository interface {
	GetByCPF(ctx context.Context, cpf string) (*domain.User, error)
}
