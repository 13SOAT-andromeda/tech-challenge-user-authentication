# Specification: User Validation Lambda

## Overview
Develop an AWS Lambda function in Go that authenticates users by validating their CPF against a PostgreSQL database (RDS) and, upon success, issues a JWT token, storing a record in DynamoDB.

## Functional Requirements
- **Request Handling**: Receive a GET request with the user's CPF in the `x-user-cpf` header.
- **CPF Validation**: Validate the CPF format using the regex: `^(\d{3}\.\d{3}\.\d{3}\-\d{2})?$`.
- **User Validation (PostgreSQL)**:
    - Query the `garagedb` database to find a user with the provided CPF.
    - Ensure the user exists and the `active` field is true.
    - Use Prepared Statements to prevent SQL injection.
- **Authentication & Token Generation**:
    - If valid and active, generate a JWT token with a 24-hour expiration.
    - Use the JWT secret from an environment variable.
- **Token Persistence (DynamoDB)**:
    - Store the token record in the `user-auth-tokens` table.
    - Schema: `token_hash` (Partition Key), `expiration_timestamp` (TTL).
- **Responses**:
    - **200 OK**: Return a JSON body with `{"token": "...", "expires_at": "..."}`.
    - **400 Bad Request**: For invalid CPF format.
    - **404 Not Found**: If the user does not exist or is inactive.

## Non-Functional Requirements
- **Architecture**: Clean Architecture / Hexagonal.
- **Language**: Go (Golang) with AWS SDK v2.
- **Runtime**: AWS Lambda (`github.com/aws/aws-lambda-go/lambda`).
- **Security**: Database credentials and JWT secret managed via environment variables.
- **Database Connection**: Use connection pools for PostgreSQL.
- **Compatibility**: Designed for LocalStack environment.

## Acceptance Criteria
- [ ] Lambda extracts CPF from header correctly.
- [ ] CPF validation correctly identifies valid and invalid formats.
- [ ] Successful database lookup returns 200 OK with token and expiration.
- [ ] Non-existent or inactive users result in 404 Not Found.
- [ ] Token is persisted in DynamoDB with correct TTL.
- [ ] Implementation follows the provided hexagonal structure.

## Out of Scope
- Full CRUD operations for users.
- AWS Secrets Manager integration (using environment variables for now).
- External API Gateway configuration (using LocalStack/CLI for testing).
