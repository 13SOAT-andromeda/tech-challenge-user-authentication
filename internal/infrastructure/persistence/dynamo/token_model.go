package dynamo

type TokenModel struct {
	PK        string `dynamodbav:"pk"`
	Token     string `dynamodbav:"token"`
	ExpiresAt int64  `dynamodbav:"expires_at"`
}
