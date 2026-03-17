package repositories

import (
	"context"
	"fmt"
	"tech-challenge-user-validation/internal/core/domain"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type GORMDynamoTokenRepository struct {
	client    DynamoDBAPI
	tableName string
}

func NewDynamoTokenRepository(client DynamoDBAPI, tableName string) *GORMDynamoTokenRepository {
	return &GORMDynamoTokenRepository{
		client:    client,
		tableName: tableName,
	}
}

func (d *GORMDynamoTokenRepository) Save(ctx context.Context, token *domain.Token) error {
	item, err := attributevalue.MarshalMap(token)
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
