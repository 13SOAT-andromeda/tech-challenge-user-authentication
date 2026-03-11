# Go Gin Lambda on LocalStack

This project is a simple Go application using the Gin framework, designed to run as an AWS Lambda function.

## Prerequisites

- Go 1.22+ (Installed in `$HOME/sdk/go` by this script)
- Docker and Docker Compose
- `awslocal` CLI

## Structure

- `go-gin-lambda/main.go`: The Go application logic.
- `go-gin-lambda/docker-compose.yml`: LocalStack configuration.
- `deploy.sh`: Script to build, zip, and deploy to LocalStack.
- `Makefile`: Alternative build/deploy tool.

## How to Run

1. **Start LocalStack:**
   ```bash
   cd go-gin-lambda
   docker-compose up -d
   cd ..
   ```

2. **Deploy the application:**
   ```bash
   ./deploy.sh
   ```

3. **Test the endpoint:**
   The `deploy.sh` script will output a `curl` command. It will look something like:
   ```bash
   curl http://localhost:4566/restapis/<API_ID>/dev/_user_request_/hello
   ```

## Development

To modify the response, edit `go-gin-lambda/main.go` and run `./deploy.sh` again.

