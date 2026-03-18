package repositories

import (
	"context"
	"tech-challenge-user-validation/internal/adapters/database/model"
	"tech-challenge-user-validation/internal/core/domain"

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
	var userModel model.UserModel
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
