# Specification: User Authentication Lambda

## Overview
Implement an AWS Lambda function in Go for user authentication validation. The service will validate the CPF format, verify the user's status in an RDS database, and generate a JWT token for active users, storing it in DynamoDB.

## Architecture
Follow the Clean Architecture / Hexagonal pattern standardized by the Andromeda group (`https://github.com/13SOAT-andromeda/tech-challenge-s1`).

### Folder Structure
- `cmd/auth/`: Application entry point and dependency injection for the Lambda.
- `internal/core/domain/`: Domain entities (e.g., User, Token).
- `internal/core/usecases/`: Business logic implementations (Auth UseCase).
- `internal/core/ports/`: Interfaces for repositories and external services.
- `internal/adapters/handlers/`: Lambda handler (API delivery).
- `internal/adapters/repositories/`: Database integration files (RDS, DynamoDB).
- `internal/infrastructure/`: AWS SDK configuration and database connection setups.

## Functional Requirements
1.  **CPF Format Validation**: Validate the CPF received in the header using the regex `^(\d{3}\.\d{3}\.\d{3}\-\d{2})?$`.
2.  **RDS User Query (Integration Only)**: 
    - Create the integration file for querying the `users` table in `garagedb`.
    - Follow the structure of the model analyzed in the reference repository.
    - Logic for the actual query should not be implemented.
3.  **User Status Validation**: Verify if the user is active (true/false).
4.  **JWT Token Generation**: 
    - Generate a JWT token with a 24-hour expiration.
    - Include a unique identifier (`jti`) in the token.
5.  **DynamoDB Token Storage (Integration Only)**:
    - Create the integration file for storing tokens in the `user-authentication-token` table.
    - Fields: `token_id` (value from JWT `jti`), `user_id` (id returned from RDS query).
    - Table Configuration: Global Secondary Index (GSI) on `user_id`, TTL enabled.
    - Logic for the actual storage should not be implemented.

## Out of Scope
- Implementation of actual query logic for RDS and DynamoDB (integration files only).
- Usage of `go-gin-lambda` adapter; use the standard AWS Lambda handler.

## Acceptance Criteria
- [ ] Folder structure follows Clean Architecture / Hexagonal pattern.
- [ ] `go-gin-lambda` folder is removed.
- [ ] CPF validation correctly uses the specified regex.
- [ ] User status validation logic is implemented in the UseCase.
- [ ] JWT generation includes `jti` and 24h expiration.
- [ ] Integration files for RDS and DynamoDB are created with correct models and table structures.
