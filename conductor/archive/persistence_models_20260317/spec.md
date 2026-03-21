# Specification: Implement Persistence Models for RDS and DynamoDB

## Overview
Implement the persistence layer models for RDS (PostgreSQL) and DynamoDB, ensuring clean naming conventions, proper tag mapping, and isolation from the core domain.

## Functional Requirements
- **PostgreSQL (RDS) Persistence Model:**
    - File: `internal/infrastructure/persistence/postgres/user_model.go`
    - Name: `UserModel`
    - Tags: Use `db:"..."` with `snake_case` mapping.
    - Fields:
        - `ID int64` (PK)
        - `Document string`
        - `IsActive bool`
        - `CreatedAt time.Time`
    - Method: `ToDomain() domain.User`
- **DynamoDB (Auth Token) Persistence Model:**
    - File: `internal/infrastructure/persistence/dynamo/token_model.go`
    - Name: `TokenModel`
    - Tags: Use `dynamodbav:"..."` with `snake_case` mapping.
    - Fields:
        - `PK string` (User CPF/Document)
        - `Token string` (JWT)
        - `ExpiresAt int64` (Unix timestamp in seconds for TTL)
- **Standardization:**
    - Remove technology-specific names from functions and factories (e.g., `NewUserRepository` instead of `NewGORMUserRepository`).
    - Use clean names like `Save` or `Client` for technology-agnostic interfaces.
    - Ensure zero leakage of database tags into the `core/domain` layer.

## Non-Functional Requirements
- **Maintainability:** Follow Clean Architecture principles by isolating persistence concerns from the domain layer.
- **Portability:** Use standard database tags and AWS SDK v2 tags where appropriate.

## Acceptance Criteria
- Both `user_model.go` and `token_model.go` are implemented in their respective packages.
- All models correctly map to their database types using provided tags.
- The `ToDomain()` method in `UserModel` correctly converts the persistence entity to a `domain.User`.
- Repository factory names are standardized to remove implementation details.
- No database tags exist in the `internal/core/domain/` package.

## Out of Scope
- Implementation of the full repository logic (only models and factory names).
- Migrations or database schema creation.
- JWT generation logic.
