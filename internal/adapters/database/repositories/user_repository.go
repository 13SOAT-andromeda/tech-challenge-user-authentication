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

func (u *userRepository) Search(ctx context.Context, params ports.UserSearch) []domain.User {
	models := []user.Model{}
	q := u.db.Model(&models)
	if params.Name != "" {
		q.Where("lower(name) LIKE ?", "%"+strings.ToLower(params.Name)+"%")
	}
	if params.Email != "" {
		q.Where("lower(email) LIKE ?", "%"+strings.ToLower(params.Email)+"%")
	}
	if params.Contact != "" {
		q.Where("lower(contact) LIKE ?", "%"+strings.ToLower(params.Contact)+"%")
	}
	q.Find(&models)

	var users []domain.User
	for _, m := range models {
		users = append(users, *m.ToDomain())
	}
	return users
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
