package repositories

import (
	"context"
	"testing"
	"time"

	"tech-challenge-user-validation/internal/adapters/database/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type mockSessionDynamoClient struct {
	putItemFunc func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	getItemFunc func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

func (m *mockSessionDynamoClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.putItemFunc(ctx, params, optFns...)
}

func (m *mockSessionDynamoClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.getItemFunc(ctx, params, optFns...)
}

func TestSessionRepository_Save(t *testing.T) {
	ctx := context.Background()
	session := model.SessionModel{
		SessionID: "550e8400-e29b-41d4-a716-446655440000",
		UserID:    "1",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	t.Run("should succeed when put item succeeds", func(t *testing.T) {
		mock := &mockSessionDynamoClient{
			putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				if *params.TableName != "user-auth-tokens" {
					t.Errorf("expected table user-auth-tokens, got %s", *params.TableName)
				}
				return &dynamodb.PutItemOutput{}, nil
			},
			getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{}, nil
			},
		}

		repo := NewSessionRepository(mock, "user-auth-tokens")
		err := repo.Save(ctx, session)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestSessionRepository_FindBySessionID(t *testing.T) {
	ctx := context.Background()
	sessionID := "550e8400-e29b-41d4-a716-446655440000"
	expiresAt := time.Now().Add(7 * 24 * time.Hour).UTC().Truncate(time.Second)

	t.Run("should return session when item exists", func(t *testing.T) {
		expected := model.SessionModel{
			SessionID: sessionID,
			UserID:    "1",
			ExpiresAt: expiresAt,
		}
		item, _ := attributevalue.MarshalMap(expected)

		mock := &mockSessionDynamoClient{
			putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				return &dynamodb.PutItemOutput{}, nil
			},
			getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				keyVal := params.Key["token_id"].(*types.AttributeValueMemberS).Value
				if keyVal != sessionID {
					t.Errorf("expected token_id=%s, got %s", sessionID, keyVal)
				}
				return &dynamodb.GetItemOutput{Item: item}, nil
			},
		}

		repo := NewSessionRepository(mock, "user-auth-tokens")
		result, err := repo.FindBySessionID(ctx, sessionID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected session, got nil")
		}
		if result.SessionID != sessionID {
			t.Errorf("expected SessionID=%s, got %s", sessionID, result.SessionID)
		}
		if result.UserID != "1" {
			t.Errorf("expected UserID=1, got %s", result.UserID)
		}
	})

	t.Run("should return nil when item does not exist", func(t *testing.T) {
		mock := &mockSessionDynamoClient{
			putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				return &dynamodb.PutItemOutput{}, nil
			},
			getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{Item: nil}, nil
			},
		}

		repo := NewSessionRepository(mock, "user-auth-tokens")
		result, err := repo.FindBySessionID(ctx, sessionID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %+v", result)
		}
	})
}
