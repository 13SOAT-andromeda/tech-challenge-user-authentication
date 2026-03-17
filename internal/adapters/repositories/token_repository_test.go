package repositories

import (
	"context"
	"testing"

	"tech-challenge-user-validation/internal/core/domain"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type mockDynamoDBAPI struct {
	putItemFunc func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func (m *mockDynamoDBAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.putItemFunc(ctx, params, optFns...)
}

func TestDynamoTokenRepository_Save(t *testing.T) {
	ctx := context.Background()
	token := &domain.Token{
		TokenID: "jti-123",
		UserID:  456,
	}

	t.Run("should succeed when put item succeeds", func(t *testing.T) {
		mockSvc := &mockDynamoDBAPI{
			putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				if *params.TableName != "user-auth-tokens" {
					t.Errorf("expected table user-auth-tokens, got %s", *params.TableName)
				}
				return &dynamodb.PutItemOutput{}, nil
			},
		}

		repo := NewDynamoTokenRepository(mockSvc, "user-auth-tokens")
		err := repo.Save(ctx, token)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
