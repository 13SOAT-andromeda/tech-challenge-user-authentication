#!/bin/bash
set -e

# Configurações do Projeto
PROJECT_DIR="go-gin-lambda"
BINARY_NAME="bootstrap"
ZIP_NAME="function.zip"
FUNCTION_NAME="go-gin-lambda"
ROLE_ARN="arn:aws:iam::000000000000:role/lambda-role"
GO_BIN="$HOME/sdk/go/bin/go"

# Configurações do Banco de Dados (Environment Variables)
DB_HOST="host.docker.internal" # No LocalStack, use este host para acessar o banco no Docker local
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="password"
DB_NAME="minha_app"

echo "Building binary..."
cd $PROJECT_DIR
GOOS=linux GOARCH=amd64 $GO_BIN build -o $BINARY_NAME main.go

echo "Zipping binary..."
zip $ZIP_NAME $BINARY_NAME

echo "Creating Lambda function in LocalStack..."
awslocal lambda delete-function --function-name $FUNCTION_NAME || true

# Adicionado o parâmetro --environment para injetar as variáveis
awslocal lambda create-function \
    --function-name $FUNCTION_NAME \
    --runtime provided.al2023 \
    --handler bootstrap \
    --role $ROLE_ARN \
    --zip-file fileb://$ZIP_NAME \
    --environment "Variables={DB_HOST=$DB_HOST,DB_PORT=$DB_PORT,DB_USER=$DB_USER,DB_PASSWORD=$DB_PASSWORD,DB_NAME=$DB_NAME}"

echo "Setting up API Gateway..."
# O deploy.sh original captura o API_ID e configura os recursos
API_ID=$(awslocal apigateway create-rest-api --name 'GinAPI' --query 'id' --output text)
PARENT_RESOURCE_ID=$(awslocal apigateway get-resources --rest-api-id $API_ID --query 'items[0].id' --output text)

RESOURCE_ID=$(awslocal apigateway create-resource --rest-api-id $API_ID --parent-id $PARENT_RESOURCE_ID --path-part 'hello' --query 'id' --output text)

awslocal apigateway put-method --rest-api-id $API_ID --resource-id $RESOURCE_ID --http-method GET --authorization-type "NONE"

awslocal apigateway put-integration \
    --rest-api-id $API_ID \
    --resource-id $RESOURCE_ID \
    --http-method GET \
    --type AWS_PROXY \
    --integration-http-method POST \
    --uri arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:000000000000:function:$FUNCTION_NAME/invocations

awslocal apigateway create-deployment --rest-api-id $API_ID --stage-name dev

echo "---------------------------------------------------"
echo "Deployment Complete!"
echo "You can test the endpoint with:"
echo "curl http://localhost:4566/restapis/$API_ID/dev/_user_request_/hello"
echo "---------------------------------------------------"