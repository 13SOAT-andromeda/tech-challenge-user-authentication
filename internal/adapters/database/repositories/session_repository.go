package repositories

import (
	"context"
	"fmt"
	"tech-challenge-user-validation/internal/adapters/database/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type SessionDynamoClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type SessionRepository struct {
	client    SessionDynamoClient
	tableName string
}

func NewSessionRepository(client SessionDynamoClient, tableName string) *SessionRepository {
	return &SessionRepository{client: client, tableName: tableName}
}

func (r *SessionRepository) Save(ctx context.Context, s model.SessionModel) error {
	item, err := attributevalue.MarshalMap(s)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put session item: %w", err)
	}

	return nil
}
