package repositories

import (
	"context"
	"errors"
	"tech-challenge-user-validation/internal/adapters/database/model/user"
	"tech-challenge-user-validation/internal/core/domain"
	"tech-challenge-user-validation/internal/core/ports"

	"gorm.io/gorm"
)

type userRepository struct {
	*BaseRepository[user.Model]
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[user.Model](db),
	}
}

func (u *userRepository) GetByDocument(ctx context.Context, document string) (*domain.User, error) {
	model := user.Model{}
	err := u.db.Where("document = ?", document).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (u *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	model := user.Model{}
	err := u.db.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}
