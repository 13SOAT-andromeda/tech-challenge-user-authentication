package repositories

import (
	"context"
	"regexp"
	"testing"
	"time"

	"tech-challenge-user-validation/internal/core/ports"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)
	defer db.DB()

	repo := NewUserRepository(db)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "contact", "document", "is_active", "password", "role", "created_at", "updated_at"}).
			AddRow(1, "Test User", "test@example.com", "123", "doc123", true, "hashed", "user", now, now)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "User" WHERE email = $1 AND "User"."deleted_at" IS NULL ORDER BY "User"."id" LIMIT $2`)).
			WithArgs("test@example.com", 1).
			WillReturnRows(rows)

		userDomain, err := repo.GetByEmail(context.Background(), "test@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, userDomain)
		assert.Equal(t, "test@example.com", userDomain.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "User" WHERE email = $1 AND "User"."deleted_at" IS NULL ORDER BY "User"."id" LIMIT $2`)).
			WithArgs("notfound@example.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{}))

		userDomain, err := repo.GetByEmail(context.Background(), "notfound@example.com")
		assert.NoError(t, err)
		assert.Nil(t, userDomain)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Search(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)
	defer db.DB()

	repo := NewUserRepository(db)
	now := time.Now()

	t.Run("search by name", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "contact", "document", "is_active", "password", "role", "created_at", "updated_at"}).
			AddRow(1, "John Doe", "john@test.com", "123", "doc123", true, "hashed", "user", now, now)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "User" WHERE lower(name) LIKE $1 AND "User"."deleted_at" IS NULL`)).
			WithArgs("%john%").
			WillReturnRows(rows)

		users := repo.Search(context.Background(), ports.UserSearch{Name: "John"})
		assert.Len(t, users, 1)
		assert.Equal(t, "John Doe", users[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_GetByDocument(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)
	defer db.DB()

	repo := NewUserRepository(db)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "contact", "document", "is_active", "password", "role", "created_at", "updated_at"}).
			AddRow(1, "Test User", "test@example.com", "123", "doc123", true, "hashed", "user", now, now)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "User" WHERE document = $1 AND "User"."deleted_at" IS NULL ORDER BY "User"."id" LIMIT $2`)).
			WithArgs("doc123", 1).
			WillReturnRows(rows)

		userDomain, err := repo.GetByDocument(context.Background(), "doc123")
		assert.NoError(t, err)
		assert.NotNil(t, userDomain)
		assert.Equal(t, "doc123", userDomain.Document)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
