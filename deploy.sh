#!/bin/bash
set -e

FUNCTION_NAME="auth-lambda"
BINARY_NAME="bootstrap"
ZIP_NAME="function.zip"
ROLE_ARN="arn:aws:iam::000000000000:role/lambda-role"

echo "Building binary..."
GOOS=linux GOARCH=amd64 go build -o $BINARY_NAME cmd/auth/main.go

echo "Zipping binary..."
zip $ZIP_NAME $BINARY_NAME

echo "Deploying to LocalStack..."
awslocal lambda delete-function --function-name $FUNCTION_NAME || true
awslocal lambda create-function \
    --function-name $FUNCTION_NAME \
    --runtime provided.al2023 \
    --handler bootstrap \
    --role $ROLE_ARN \
    --zip-file fileb://$ZIP_NAME

echo "Setting up API Gateway..."
API_ID=$(awslocal apigateway create-rest-api --name 'AuthAPI' --query 'id' --output text)
PARENT_RESOURCE_ID=$(awslocal apigateway get-resources --rest-api-id $API_ID --query 'items[0].id' --output text)

# Create /auth resource
RESOURCE_ID=$(awslocal apigateway create-resource --rest-api-id $API_ID --parent-id $PARENT_RESOURCE_ID --path-part 'auth' --query 'id' --output text)

# Create GET method (since we receive CPF in header)
awslocal apigateway put-method --rest-api-id $API_ID --resource-id $RESOURCE_ID --http-method GET --authorization-type "NONE"

# Integrate with Lambda
awslocal apigateway put-integration \
    --rest-api-id $API_ID \
    --resource-id $RESOURCE_ID \
    --http-method GET \
    --type AWS_PROXY \
    --integration-http-method POST \
    --uri arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:000000000000:function:$FUNCTION_NAME/invocations

# Deploy API
awslocal apigateway create-deployment --rest-api-id $API_ID --stage-name dev

echo "---------------------------------------------------"
echo "Deployment Complete!"
echo "Test endpoint:"
echo "curl -H 'x-cpf: 123.456.789-00' http://localhost:4566/restapis/$API_ID/dev/_user_request_/auth"
echo "---------------------------------------------------"
