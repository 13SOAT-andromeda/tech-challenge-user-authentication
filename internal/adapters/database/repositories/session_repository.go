package repositories

import (
	"context"
	"fmt"
	"tech-challenge-user-validation/internal/adapters/database/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SessionDynamoClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
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

func (r *SessionRepository) FindBySessionID(ctx context.Context, sessionID string) (*model.SessionModel, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"token_id": &types.AttributeValueMemberS{Value: sessionID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get session item: %w", err)
	}
	if out.Item == nil {
		return nil, nil
	}

	var s model.SessionModel
	if err := attributevalue.UnmarshalMap(out.Item, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}
	return &s, nil
}

