package repositories

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository_GetByDocument(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbMock.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dbMock,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm db", err)
	}

	repo := NewUserRepository(gormDB)
	ctx := context.Background()
	document := "123.456.789-00"

	t.Run("should return user when found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "document", "is_active", "created_at"}).
			AddRow(1, document, true, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE document = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(document, 1).
			WillReturnRows(rows)

		user, err := repo.GetByDocument(ctx, document)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if user == nil {
			t.Fatal("expected user to be found, got nil")
		}
		if user.Document != document {
			t.Errorf("expected document %s, got %s", document, user.Document)
		}
	})

	t.Run("should return nil when user not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE document = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(document, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		user, err := repo.GetByDocument(ctx, document)

		if err != nil && err != gorm.ErrRecordNotFound {
			t.Errorf("expected no error or record not found, got %s", err)
		}
		if user != nil {
			t.Errorf("expected user to be nil, got %v", user)
		}
	})
}
