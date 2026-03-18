package repositories

import (
	"context"
	"tech-challenge-user-validation/internal/core/domain"
	"tech-challenge-user-validation/internal/infrastructure/persistence/postgres"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByDocument(ctx context.Context, document string) (*domain.User, error) {
	var userModel postgres.UserModel
	err := r.db.WithContext(ctx).Where("document = ?", document).First(&userModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	user := userModel.ToDomain()
	return &user, nil
}
