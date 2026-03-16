BINARY_NAME=bootstrap
ZIP_NAME=function.zip
FUNCTION_NAME=auth-lambda
ROLE_ARN=arn:aws:iam::000000000000:role/lambda-role

build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) cmd/auth/main.go

zip: build
	zip $(ZIP_NAME) $(BINARY_NAME)

deploy: zip
	awslocal lambda delete-function --function-name $(FUNCTION_NAME) || true
	awslocal lambda create-function \
		--function-name $(FUNCTION_NAME) \
		--runtime provided.al2023 \
		--handler bootstrap \
		--role $(ROLE_ARN) \
		--zip-file fileb://$(ZIP_NAME)

clean:
	rm -f $(BINARY_NAME) $(ZIP_NAME)

redeploy: clean deploy
