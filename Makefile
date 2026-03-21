# ─── Config ───────────────────────────────────────────────────────────────────
BINARY        := bootstrap
ZIP           := function.zip
FUNCTION_NAME := user-auth-function
RUNTIME       := provided.al2023
REGION        := us-east-1
ROLE_ARN      := arn:aws:iam::000000000000:role/lambda-role
API_NAME      := user-auth-api
API_STAGE     := local

# Run awslocal inside the LocalStack container to avoid WSL2 host→container
# networking issues. The zip is copied into the container before deploy.
AWSLOCAL = docker exec localstack awslocal

# ─── Environment Variables ─────────────────────────────────────────────────────
# Override any of these on the command line: make deploy JWT_SECRET=mysecret
JWT_SECRET    ?= local-dev-secret
DB_HOST       ?= postgres
DB_USER       ?= postgres
DB_PASSWORD   ?= postgres
DB_NAME       ?= authdb
DB_PORT       ?= 5432
DYNAMO_TABLE  ?= user-auth-tokens

# Login parameters for make curl / make invoke
DOCUMENT      ?= 123.456.789-00
PASSWORD      ?= Admin123!

ENV_VARS = Variables="{JWT_SECRET=$(JWT_SECRET),DB_HOST=$(DB_HOST),DB_USER=$(DB_USER),DB_PASSWORD=$(DB_PASSWORD),DB_NAME=$(DB_NAME),DB_PORT=$(DB_PORT),DYNAMODB_TABLE_NAME=$(DYNAMO_TABLE),AWS_ENDPOINT_URL=http://localstack:4566,AWS_ACCESS_KEY_ID=test,AWS_SECRET_ACCESS_KEY=test,AWS_REGION=$(REGION)}"

# ─── Build ─────────────────────────────────────────────────────────────────────
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY) ./cmd/main.go
	@echo "Build complete: $(BINARY)"

.PHONY: zip
zip: build
	zip -j $(ZIP) $(BINARY)
	@echo "Package ready: $(ZIP)"

# ─── LocalStack ────────────────────────────────────────────────────────────────
.PHONY: localstack-up
localstack-up:
	docker compose up -d
	@echo "Waiting for LocalStack..."
	@until docker exec localstack curl -sf http://localhost:4566/_localstack/health 2>/dev/null | grep -q '"lambda": "available"'; do \
		printf '.'; sleep 2; \
	done
	@echo ""
	@echo "Waiting for PostgreSQL..."
	@until docker exec localstack-postgres pg_isready -U postgres > /dev/null 2>&1; do \
		printf '.'; sleep 2; \
	done
	@echo ""
	@echo "Services are ready."

.PHONY: localstack-down
localstack-down:
	docker compose down
	@echo "Services stopped."

.PHONY: localstack-clean
localstack-clean:
	docker compose down -v
	rm -rf .localstack
	@echo "All LocalStack data removed."

# ─── Infrastructure ────────────────────────────────────────────────────────────
.PHONY: setup-infra
setup-infra: setup-aws setup-api-gateway
	@echo "Infrastructure ready."

.PHONY: setup-aws
setup-aws:
	@echo "Creating IAM role..."
	$(AWSLOCAL) iam create-role \
		--role-name lambda-role \
		--assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}' \
		--region $(REGION) 2>/dev/null || true

	@echo "Creating DynamoDB table..."
	$(AWSLOCAL) dynamodb create-table \
		--table-name $(DYNAMO_TABLE) \
		--attribute-definitions AttributeName=token_id,AttributeType=S \
		--key-schema AttributeName=token_id,KeyType=HASH \
		--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
		--region $(REGION) 2>/dev/null || true

.PHONY: setup-api-gateway
setup-api-gateway:
	@EXISTING=$$($(AWSLOCAL) apigateway get-rest-apis \
		--region $(REGION) \
		--query 'items[?name==`$(API_NAME)`].id' --output text | awk '{print $$1}'); \
	if [ -n "$$EXISTING" ]; then \
		echo "API Gateway already exists: $$EXISTING"; \
		echo "Endpoint: http://localhost:4566/restapis/$$EXISTING/$(API_STAGE)/_user_request_/sessions"; \
	else \
		echo "Creating API Gateway..."; \
		API_ID=$$($(AWSLOCAL) apigateway create-rest-api \
			--name $(API_NAME) \
			--region $(REGION) \
			--query 'id' --output text); \
		ROOT_ID=$$($(AWSLOCAL) apigateway get-resources \
			--rest-api-id $$API_ID \
			--region $(REGION) \
			--query 'items[?path==`/`].id' --output text); \
		RESOURCE_ID=$$($(AWSLOCAL) apigateway create-resource \
			--rest-api-id $$API_ID \
			--parent-id $$ROOT_ID \
			--path-part sessions \
			--region $(REGION) \
			--query 'id' --output text); \
		$(AWSLOCAL) apigateway put-method \
			--rest-api-id $$API_ID \
			--resource-id $$RESOURCE_ID \
			--http-method POST \
			--authorization-type NONE \
			--region $(REGION) > /dev/null; \
		$(AWSLOCAL) apigateway put-integration \
			--rest-api-id $$API_ID \
			--resource-id $$RESOURCE_ID \
			--http-method POST \
			--type AWS_PROXY \
			--integration-http-method POST \
			--uri "arn:aws:apigateway:$(REGION):lambda:path/2015-03-31/functions/arn:aws:lambda:$(REGION):000000000000:function:$(FUNCTION_NAME)/invocations" \
			--region $(REGION) > /dev/null; \
		$(AWSLOCAL) apigateway create-deployment \
			--rest-api-id $$API_ID \
			--stage-name $(API_STAGE) \
			--region $(REGION) > /dev/null; \
		echo "API Gateway ready."; \
		echo "Endpoint: http://localhost:4566/restapis/$$API_ID/$(API_STAGE)/_user_request_/sessions"; \
	fi

# ─── Deploy ────────────────────────────────────────────────────────────────────
.PHONY: deploy
deploy: zip
	@echo "Copying zip into LocalStack container..."
	docker cp $(ZIP) localstack:/tmp/$(ZIP)

	@echo "Deploying $(FUNCTION_NAME)..."
	@if $(AWSLOCAL) lambda get-function --function-name $(FUNCTION_NAME) --region $(REGION) > /dev/null 2>&1; then \
		echo "Updating function code..."; \
		$(AWSLOCAL) lambda update-function-code \
			--function-name $(FUNCTION_NAME) \
			--zip-file fileb:///tmp/$(ZIP) \
			--region $(REGION); \
		echo "Updating environment variables..."; \
		$(AWSLOCAL) lambda update-function-configuration \
			--function-name $(FUNCTION_NAME) \
			--environment "$(ENV_VARS)" \
			--region $(REGION); \
	else \
		echo "Creating function..."; \
		$(AWSLOCAL) lambda create-function \
			--function-name $(FUNCTION_NAME) \
			--runtime $(RUNTIME) \
			--handler $(BINARY) \
			--role $(ROLE_ARN) \
			--zip-file fileb:///tmp/$(ZIP) \
			--environment "$(ENV_VARS)" \
			--region $(REGION); \
	fi
	@echo "Deploy complete."

.PHONY: redeploy
redeploy: zip
	@echo "Redeploying code only..."
	docker cp $(ZIP) localstack:/tmp/$(ZIP)
	$(AWSLOCAL) lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--zip-file fileb:///tmp/$(ZIP) \
		--region $(REGION)
	@echo "Redeploy complete."

# ─── Seed ──────────────────────────────────────────────────────────────────────
.PHONY: seed
seed:
	@echo "Seeding database..."
	docker exec -i localstack-postgres psql -U $(DB_USER) -d $(DB_NAME) < scripts/seed.sql
	@echo "Seed complete. User: 123.456.789-00 / Admin123!"

# ─── Test ──────────────────────────────────────────────────────────────────────
.PHONY: curl
curl:
	@API_ID=$$($(AWSLOCAL) apigateway get-rest-apis \
		--region $(REGION) \
		--query 'items[?name==`$(API_NAME)`].id' --output text | awk '{print $$1}'); \
	docker exec localstack curl --silent --location --request POST \
		"http://localhost:4566/restapis/$$API_ID/$(API_STAGE)/_user_request_/sessions" \
		--header 'Content-Type: application/json' \
		--data-raw '{"document":"$(DOCUMENT)","password":"$(PASSWORD)"}' | cat

.PHONY: invoke
invoke:
	$(AWSLOCAL) lambda invoke \
		--function-name $(FUNCTION_NAME) \
		--payload '{"httpMethod":"POST","path":"/sessions","body":"{\"document\":\"$(DOCUMENT)\",\"password\":\"$(PASSWORD)\"}"}' \
		--region $(REGION) \
		/tmp/response.json
	docker exec localstack cat /tmp/response.json

.PHONY: logs
logs:
	$(AWSLOCAL) lambda get-function \
		--function-name $(FUNCTION_NAME) \
		--region $(REGION)

# ─── Shortcuts ─────────────────────────────────────────────────────────────────
.PHONY: local
local: localstack-up setup-infra deploy

.PHONY: clean
clean:
	rm -f $(BINARY) $(ZIP)

.PHONY: test
test:
	go test ./... -v

.PHONY: help
help:
	@echo ""
	@echo "  make local            Start LocalStack + deploy everything"
	@echo "  make localstack-up    Start docker services"
	@echo "  make localstack-down  Stop docker services"
	@echo "  make localstack-clean Stop and remove all data"
	@echo "  make setup-infra      Create IAM role + DynamoDB table + API Gateway"
	@echo "  make build            Compile binary for Linux/amd64"
	@echo "  make zip              Build and zip the binary"
	@echo "  make deploy           Build, zip and deploy to LocalStack"
	@echo "  make redeploy         Re-zip and push code only (faster)"
	@echo "  make seed             Insert test user into Postgres"
	@echo "  make curl             POST /sessions via API Gateway"
	@echo "  make invoke           Direct Lambda invoke (no API Gateway)"
	@echo "  make logs             Show function metadata"
	@echo "  make clean            Remove build artifacts"
	@echo "  make test             Run all tests"
	@echo ""
	@echo "  Override credentials: make curl DOCUMENT=123.456.789-00 PASSWORD=pass"
	@echo "  Override env vars:    make deploy JWT_SECRET=mysecret DB_PASSWORD=pass"
	@echo ""
