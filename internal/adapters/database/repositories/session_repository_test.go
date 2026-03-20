package repositories

import (
	"context"
	"tech-challenge-user-validation/internal/adapters/database/model"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type mockDynamoDBAPI struct {
	putItemFunc func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func (m *mockDynamoDBAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.putItemFunc(ctx, params, optFns...)
}

func TestTokenRepository_Save(t *testing.T) {
	ctx := context.Background()
	pk := "12345678900"
	token := "some-jwt-token"
	expiresAt := int64(1679062200)

	t.Run("should succeed when put item succeeds", func(t *testing.T) {
		mockSvc := &mockDynamoDBAPI{
			putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				if *params.TableName != "user-auth-tokens" {
					t.Errorf("expected table user-auth-tokens, got %s", *params.TableName)
				}
				return &dynamodb.PutItemOutput{}, nil
			},
		}

		repo := model.NewTokenRepository(mockSvc, "user-auth-tokens")
		err := repo.Save(ctx, pk, token, expiresAt)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
