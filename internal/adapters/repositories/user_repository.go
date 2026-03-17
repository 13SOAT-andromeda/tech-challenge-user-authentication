package repositories

import (
	"context"
	"tech-challenge-user-validation/internal/core/domain"

	"gorm.io/gorm"
)

type GORMUserRepository struct {
	db *gorm.DB
}

func NewGORMUserRepository(db *gorm.DB) *GORMUserRepository {
	return &GORMUserRepository{
		db: db,
	}
}

func (r *GORMUserRepository) GetByDocument(ctx context.Context, Document string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("document = ?", Document).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
