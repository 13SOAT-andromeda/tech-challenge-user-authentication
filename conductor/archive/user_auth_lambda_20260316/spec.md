# Specification: User Auth Lambda (Go/GORM/DynamoDB)

## Overview
Develop an AWS Lambda function in Go to authenticate users by validating their CPF against a PostgreSQL (RDS) database and persisting a JWT session in DynamoDB.

## Functional Requirements
- **Extraction**: Retrieve the CPF from the `x-user-cpf` request header.
- **CPF Validation**: Ensure the CPF matches the regex: `^(\d{3}\.\d{3}\.\d{3}\-\d{2})?$`.
- **RDS User Query (PostgreSQL)**:
    - Target database: `garagedb`.
    - Implement a secure `GetByDocument` method using prepared statements.
    - Verify if the user exists and the `active` field is true.
- **Authentication & Persistence**:
    - On success, generate a JWT token (expires in 24h) using `golang-jwt/jwt/v5`.
    - Persist the token record in the `user-auth-tokens` DynamoDB table.
- **Responses**:
    - **200 OK**: Return a JSON body with the token on success.
    - **400 Bad Request**: Return a JSON body `{"error": "invalid format"}` for invalid CPF.
    - **404 Not Found**: Return a JSON body `{"error": "user not found or inactive"}` if the user is missing or disabled.

## Non-Functional Requirements
- **Architecture**: Clean Architecture / Hexagonal, based on `13SOAT-andromeda/tech-challenge-s1`.
- **Language**: Go (Golang) using AWS SDK v2.
- **Database Driver**: `jackc/pgx/v5` for PostgreSQL.
- **Framework Constraint**: Use the standard `aws-lambda-go` handler; do not use `gin-lambda`.
- **Infrastructure**: Include initialization scripts for LocalStack (RDS/DynamoDB).

## Acceptance Criteria
- [ ] Lambda correctly extracts CPF from the `x-user-cpf` header.
- [ ] Invalid CPF formats return a 400 response.
- [ ] Non-existent or inactive users in RDS return a 404 response.
- [ ] Active users generate a 200 response with a valid JWT.
- [ ] Tokens are successfully saved to DynamoDB.
- [ ] Database credentials and host are managed through environment variables.
- [ ] Clean Architecture structure is maintained (Handler, Usecase, Repository layers).

## Out of Scope
- Full CRUD for users.
- API Gateway configuration beyond what's needed for the Lambda.
- Advanced monitoring or logging.
