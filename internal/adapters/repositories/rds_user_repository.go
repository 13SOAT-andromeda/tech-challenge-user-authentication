package repositories

import (
	"context"
	"tech-challenge-user-validation/internal/core/domain"
)

type RDSUserRepository struct {
	// DB connection would go here
}

func NewRDSUserRepository() *RDSUserRepository {
	return &RDSUserRepository{}
}

func (r *RDSUserRepository) GetByCPF(ctx context.Context, cpf string) (*domain.User, error) {
	// Logic not implemented as per spec
	return nil, nil
}
