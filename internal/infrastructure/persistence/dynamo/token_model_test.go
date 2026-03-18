package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestTokenModel_Marshaling(t *testing.T) {
	tokenModel := TokenModel{
		PK:        "12345678900",
		Token:     "some-jwt-token",
		ExpiresAt: 1679062200,
	}

	av, err := attributevalue.MarshalMap(tokenModel)
	assert.NoError(t, err)

	assert.Equal(t, &types.AttributeValueMemberS{Value: "12345678900"}, av["pk"])
	assert.Equal(t, &types.AttributeValueMemberS{Value: "some-jwt-token"}, av["token"])
	assert.Equal(t, &types.AttributeValueMemberN{Value: "1679062200"}, av["expires_at"])
}
