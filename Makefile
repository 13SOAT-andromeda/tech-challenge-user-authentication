GO_BIN ?= $(HOME)/sdk/go/bin/go
BINARY_NAME=bootstrap
ZIP_NAME=function.zip
FUNCTION_NAME=go-gin-lambda
ROLE_ARN=arn:aws:iam::000000000000:role/lambda-role

build:
	cd go-gin-lambda && GOOS=linux GOARCH=amd64 $(GO_BIN) build -o $(BINARY_NAME) main.go

zip: build
	cd go-gin-lambda && zip $(ZIP_NAME) $(BINARY_NAME)

deploy: zip
	awslocal lambda create-function 
		--function-name $(FUNCTION_NAME) 
		--runtime provided.al2023 
		--handler bootstrap 
		--role $(ROLE_ARN) 
		--zip-file fileb://go-gin-lambda/$(ZIP_NAME)

create-api:
	awslocal apigateway create-rest-api --name 'GinAPI'
	# Note: In a real script, you'd capture the API ID and parent resource ID.
	# This Makefile is a template. I'll provide a shell script for full automation.

clean:
	rm -f go-gin-lambda/$(BINARY_NAME) go-gin-lambda/$(ZIP_NAME)

redeploy: clean deploy
