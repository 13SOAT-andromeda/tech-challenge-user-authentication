package repositories

import (
	"context"
	"errors"
	"strings"
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

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	model := user.Model{}
	err := u.db.Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (u *userRepository) Search(ctx context.Context, params ports.UserSearch) []domain.User {
	var models []user.Model
	query := u.db.Model(&user.Model{})
	if params.Name != "" {
		query = query.Where("lower(name) LIKE ?", "%"+strings.ToLower(params.Name)+"%")
	}
	query.Find(&models)
	users := make([]domain.User, 0, len(models))
	for _, m := range models {
		users = append(users, *m.ToDomain())
	}
	return users
}
