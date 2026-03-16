package ports

import (
	"context"
	"tech-challenge-user-validation/internal/core/domain"
)

type TokenRepository interface {
	Save(ctx context.Context, token *domain.Token) error
}
