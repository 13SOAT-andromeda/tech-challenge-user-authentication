package repositories

import (
	"context"
	"fmt"
	"tech-challenge-user-validation/internal/adapters/database/model"
	"tech-challenge-user-validation/internal/core/ports"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type SessionRepository struct {
	client    DynamoClient
	tableName string
}

func NewSessionRepository(client DynamoClient, tableName string) *SessionRepository {
	return &SessionRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *SessionRepository) Create(ctx context.Context, sessionID string, userID string, expiresAt int64) (*ports.Session, error) {
	sessionModel := model.SessionModel{
		PK:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	item, err := attributevalue.MarshalMap(sessionModel)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put session item: %w", err)
	}

	return sessionModel.ToDomain(), nil
}

func (r *SessionRepository) GetByID(ctx context.Context, sessionID string) (*ports.Session, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"pk": sessionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}

	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get session item: %w", err)
	}

	if out.Item == nil {
		return nil, nil
	}

	var sessionModel model.SessionModel
	err = attributevalue.UnmarshalMap(out.Item, &sessionModel)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return sessionModel.ToDomain(), nil
}
