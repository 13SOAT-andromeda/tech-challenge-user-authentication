package repositories

import (
	"context"
	"regexp"
	"testing"
	"time"

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

func TestUserRepository_GetByDocument(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)
	defer db.DB()

	repo := NewUserRepository(db)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",
			"password", "role", "person_id",
			"Person__id", "Person__created_at", "Person__updated_at", "Person__deleted_at",
			"Person__name", "Person__email", "Person__contact", "Person__document",
			"Person__is_active",
		}).AddRow(
			1, now, now, nil,
			"hashed", "user", 2,
			2, now, now, nil,
			"Test User", "test@example.com", "123", "doc123",
			true,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "User"."id","User"."created_at","User"."updated_at","User"."deleted_at","User"."password","User"."role","User"."person_id","Person"."id" AS "Person__id","Person"."created_at" AS "Person__created_at","Person"."updated_at" AS "Person__updated_at","Person"."deleted_at" AS "Person__deleted_at","Person"."name" AS "Person__name","Person"."email" AS "Person__email","Person"."contact" AS "Person__contact","Person"."document" AS "Person__document","Person"."is_active" AS "Person__is_active","Person"."street" AS "Person__street","Person"."number" AS "Person__number","Person"."complement" AS "Person__complement","Person"."city" AS "Person__city","Person"."state" AS "Person__state","Person"."zip_code" AS "Person__zip_code" FROM "User" LEFT JOIN "Person" "Person" ON "User"."person_id" = "Person"."id" AND "Person"."deleted_at" IS NULL WHERE "Person".document = $1 AND "User"."deleted_at" IS NULL ORDER BY "User"."id" LIMIT $2`)).
			WithArgs("doc123", 1).
			WillReturnRows(rows)

		userDomain, err := repo.GetByDocument(context.Background(), "doc123")
		assert.NoError(t, err)
		assert.NotNil(t, userDomain)
		assert.NotNil(t, userDomain.Person)
		assert.Equal(t, "doc123", userDomain.Person.Document)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
