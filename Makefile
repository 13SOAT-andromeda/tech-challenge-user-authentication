include .env
export

# ─── Config ───────────────────────────────────────────────────────────────────
RUNTIME       := provided.al2023
API_NAME      := user-auth-api
API_STAGE     := local

# Run awslocal inside the LocalStack container to avoid WSL2 host→container
# networking issues. The zip is copied into the container before deploy.
AWSLOCAL = docker exec localstack awslocal

# Login parameters for make curl / make invoke
DOCUMENT      ?= 123.456.789-00
PASSWORD      ?= Admin123!

ENV_VARS = Variables="{JWT_SECRET=$(JWT_SECRET),DB_HOST=$(DB_HOST),DB_USER=$(DB_USER),DB_PASSWORD=$(DB_PASSWORD),DB_NAME=$(DB_NAME),DB_PORT=$(DB_PORT),DYNAMODB_TABLE_NAME=$(DYNAMODB_TABLE_NAME),AWS_ENDPOINT_URL=http://localstack:4566,AWS_ACCESS_KEY_ID=test,AWS_SECRET_ACCESS_KEY=test,AWS_REGION=$(AWS_REGION)}"

# ─── Build ─────────────────────────────────────────────────────────────────────
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(FUNCTION_BINARY_NAME) ./cmd/main.go
	@echo "Build complete: $(FUNCTION_BINARY_NAME)"

.PHONY: zip
zip: build
	zip -j $(FUNCTION_ZIP_NAME) $(FUNCTION_BINARY_NAME)
	@echo "Package ready: $(FUNCTION_ZIP_NAME)"

# ─── LocalStack ────────────────────────────────────────────────────────────────
.PHONY: localstack-up
localstack-up:
	docker compose up -d
	@echo "Waiting for LocalStack..."
	@until docker exec localstack curl -sf http://localhost:4566/_localstack/health 2>/dev/null | grep -q '"lambda": "available"'; do \
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
		--region $(AWS_REGION) 2>/dev/null || true

	@echo "Creating DynamoDB table..."
	$(AWSLOCAL) dynamodb create-table \
		--table-name $(DYNAMODB_TABLE_NAME) \
		--attribute-definitions AttributeName=token_id,AttributeType=S \
		--key-schema AttributeName=token_id,KeyType=HASH \
		--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
		--region $(AWS_REGION) 2>/dev/null || true

.PHONY: setup-api-gateway
setup-api-gateway:
	@EXISTING=$$($(AWSLOCAL) apigateway get-rest-apis \
		--region $(AWS_REGION) \
		--query 'items[?name==`$(API_NAME)`].id' --output text | awk '{print $$1}'); \
	if [ -n "$$EXISTING" ]; then \
		echo "API Gateway already exists: $$EXISTING"; \
		echo "Endpoint: http://localhost:4566/restapis/$$EXISTING/$(API_STAGE)/_user_request_/sessions"; \
	else \
		echo "Creating API Gateway..."; \
		API_ID=$$($(AWSLOCAL) apigateway create-rest-api \
			--name $(API_NAME) \
			--region $(AWS_REGION) \
			--query 'id' --output text); \
		ROOT_ID=$$($(AWSLOCAL) apigateway get-resources \
			--rest-api-id $$API_ID \
			--region $(AWS_REGION) \
			--query 'items[?path==`/`].id' --output text); \
		RESOURCE_ID=$$($(AWSLOCAL) apigateway create-resource \
			--rest-api-id $$API_ID \
			--parent-id $$ROOT_ID \
			--path-part sessions \
			--region $(AWS_REGION) \
			--query 'id' --output text); \
		$(AWSLOCAL) apigateway put-method \
			--rest-api-id $$API_ID \
			--resource-id $$RESOURCE_ID \
			--http-method POST \
			--authorization-type NONE \
			--region $(AWS_REGION) > /dev/null; \
		$(AWSLOCAL) apigateway put-integration \
			--rest-api-id $$API_ID \
			--resource-id $$RESOURCE_ID \
			--http-method POST \
			--type AWS_PROXY \
			--integration-http-method POST \
			--uri "arn:aws:apigateway:$(AWS_REGION):lambda:path/2015-03-31/functions/arn:aws:lambda:$(AWS_REGION):000000000000:function:$(FUNCTION_NAME)/invocations" \
			--region $(AWS_REGION) > /dev/null; \
		$(AWSLOCAL) apigateway create-deployment \
			--rest-api-id $$API_ID \
			--stage-name $(API_STAGE) \
			--region $(AWS_REGION) > /dev/null; \
		echo "API Gateway ready."; \
		echo "Endpoint: http://localhost:4566/restapis/$$API_ID/$(API_STAGE)/_user_request_/sessions"; \
	fi

# ─── Deploy ────────────────────────────────────────────────────────────────────
.PHONY: deploy
deploy: zip
	@echo "Copying zip into LocalStack container..."
	docker cp $(FUNCTION_ZIP_NAME) localstack:/tmp/$(FUNCTION_ZIP_NAME)

	@echo "Deploying $(FUNCTION_NAME)..."
	@if $(AWSLOCAL) lambda get-function --function-name $(FUNCTION_NAME) --region $(AWS_REGION) > /dev/null 2>&1; then \
		echo "Updating function code..."; \
		$(AWSLOCAL) lambda update-function-code \
			--function-name $(FUNCTION_NAME) \
			--zip-file fileb:///tmp/$(FUNCTION_ZIP_NAME) \
			--region $(AWS_REGION); \
		echo "Updating environment variables..."; \
		$(AWSLOCAL) lambda update-function-configuration \
			--function-name $(FUNCTION_NAME) \
			--environment "$(ENV_VARS)" \
			--region $(AWS_REGION); \
	else \
		echo "Creating function..."; \
		$(AWSLOCAL) lambda create-function \
			--function-name $(FUNCTION_NAME) \
			--runtime $(RUNTIME) \
			--handler $(FUNCTION_BINARY_NAME) \
			--role $(AWS_ROLE_ARN) \
			--zip-file fileb:///tmp/$(FUNCTION_ZIP_NAME) \
			--environment "$(ENV_VARS)" \
			--region $(AWS_REGION); \
	fi
	@echo "Deploy complete."

.PHONY: redeploy
redeploy: zip
	@echo "Redeploying code only..."
	docker cp $(FUNCTION_ZIP_NAME) localstack:/tmp/$(FUNCTION_ZIP_NAME)
	$(AWSLOCAL) lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--zip-file fileb:///tmp/$(FUNCTION_ZIP_NAME) \
		--region $(AWS_REGION)
	@echo "Redeploy complete."

# ─── Seed ──────────────────────────────────────────────────────────────────────
.PHONY: seed
seed:
	@echo "Seeding database..."
	docker exec -i localstack-postgres psql -U $(DB_USER) -d $(DB_NAME) -v ON_ERROR_STOP=1 < scripts/seed.sql
	@echo "Seed complete. Password for all users: Admin123!"
	@echo "  customer@example.com      (document: 11122233344)"
	@echo "  attendant@example.com     (document: 22233344455)"
	@echo "  mechanic@example.com      (document: 33344455566)"
	@echo "  administrator@example.com (document: 44455566677)"

# ─── Test ──────────────────────────────────────────────────────────────────────
.PHONY: curl
curl:
	@API_ID=$$($(AWSLOCAL) apigateway get-rest-apis \
		--region $(AWS_REGION) \
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
		--region $(AWS_REGION) \
		/tmp/response.json
	docker exec localstack cat /tmp/response.json

.PHONY: logs
logs:
	$(AWSLOCAL) lambda get-function \
		--function-name $(FUNCTION_NAME) \
		--region $(AWS_REGION)

.PHONY: sessions
sessions:
	$(AWSLOCAL) dynamodb scan \
		--table-name $(DYNAMODB_TABLE_NAME) \
		--region $(AWS_REGION)

.PHONY: sessions-flush
sessions-flush:
	@echo "Flushing all sessions from $(DYNAMODB_TABLE_NAME)..."
	@$(AWSLOCAL) dynamodb scan \
		--table-name $(DYNAMODB_TABLE_NAME) \
		--region $(AWS_REGION) \
		--query 'Items[*].token_id.S' \
		--output text | tr '\t' '\n' | while read id; do \
			$(AWSLOCAL) dynamodb delete-item \
				--table-name $(DYNAMODB_TABLE_NAME) \
				--key "{\"token_id\":{\"S\":\"$$id\"}}" \
				--region $(AWS_REGION); \
			echo "Deleted: $$id"; \
		done
	@echo "Done."

# ─── SAM Local ─────────────────────────────────────────────────────────────────

.PHONY: sam-env
sam-env:
	@echo '{"UserAuthFunction":{' \
		'"JWT_SECRET":"$(JWT_SECRET)",' \
		'"JWT_REFRESH_SECRET":"$(JWT_REFRESH_SECRET)",' \
		'"DB_HOST":"$(DB_HOST)",' \
		'"DB_USER":"$(DB_USER)",' \
		'"DB_PASSWORD":"$(DB_PASSWORD)",' \
		'"DB_NAME":"$(DB_NAME)",' \
		'"DB_PORT":"$(DB_PORT)",' \
		'"DYNAMODB_TABLE_NAME":"$(DYNAMODB_TABLE_NAME)",' \
		'"DD_TRACE_ENABLED":"$(DD_TRACE_ENABLED)",' \
		'"DD_API_KEY":"$(DD_API_KEY)",' \
		'"DD_SITE":"$(DD_SITE)",' \
		'"AWS_ENDPOINT_URL":"$(AWS_ENDPOINT_URL)",' \
		'"AWS_ACCESS_KEY_ID":"$(AWS_ACCESS_KEY_ID)",' \
		'"AWS_SECRET_ACCESS_KEY":"$(AWS_SECRET_ACCESS_KEY)",' \
		'"AWS_REGION":"$(AWS_REGION)"' \
		'}}' > sam-env.json

.PHONY: sam
sam: sam-env
	sam build
	sam local start-api \
		--port $(SAM_PORT) \
		--docker-network administrative-api \
		--env-vars sam-env.json \
		--debug \
		--warm-containers EAGER

# ─── Shortcuts ─────────────────────────────────────────────────────────────────
.PHONY: local
local: localstack-up setup-infra deploy

.PHONY: clean
clean:
	rm -f $(FUNCTION_BINARY_NAME) $(FUNCTION_ZIP_NAME) sam-env.json

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
	@echo "  make seed             Insert test users into Postgres"
	@echo "  make curl             POST /sessions via API Gateway"
	@echo "  make invoke           Direct Lambda invoke (no API Gateway)"
	@echo "  make logs             Show function metadata"
	@echo "  make sessions         List all sessions in DynamoDB"
	@echo "  make sessions-flush   Delete all sessions from DynamoDB"
	@echo "  make sam              Start SAM local API (port $(SAM_PORT))"
	@echo "  make clean            Remove build artifacts"
	@echo "  make test             Run all tests"
	@echo ""
	@echo "  Override credentials: make curl DOCUMENT=123.456.789-00 PASSWORD=pass"
	@echo "  Override env vars:    make deploy JWT_SECRET=mysecret DB_PASSWORD=pass"
	@echo ""
