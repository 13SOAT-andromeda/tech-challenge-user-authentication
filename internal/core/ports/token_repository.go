package ports

import (
	"context"
)

type TokenRepository interface {
	Save(ctx context.Context, pk string, token string, expiresAt int64) error
}
