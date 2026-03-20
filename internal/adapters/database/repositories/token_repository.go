package repositories

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TokenDynamoClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type tokenRepository struct {
	client    TokenDynamoClient
	tableName string
}

func NewTokenRepository(client TokenDynamoClient, tableName string) *tokenRepository {
	return &tokenRepository{client: client, tableName: tableName}
}

func (r *tokenRepository) Save(ctx context.Context, pk string, token string, expiresAt int64) error {
	item, err := attributevalue.MarshalMap(map[string]interface{}{
		"token_id":   pk,
		"token":      token,
		"expires_at": expiresAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put token item: %w", err)
	}

	return nil
}
