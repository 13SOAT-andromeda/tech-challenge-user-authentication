package repositories

import (
	"context"
	"fmt"
	"tech-challenge-user-validation/internal/infrastructure/persistence/dynamo"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type TokenRepository struct {
	client    Client
	tableName string
}

func NewTokenRepository(client Client, tableName string) *TokenRepository {
	return &TokenRepository{
		client:    client,
		tableName: tableName,
	}
}

func (d *TokenRepository) Save(ctx context.Context, pk string, token string, expiresAt int64) error {
	tokenModel := dynamo.TokenModel{
		PK:        pk,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	item, err := attributevalue.MarshalMap(tokenModel)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}
