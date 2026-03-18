# Implementation Plan: Implement Persistence Models for RDS and DynamoDB

## Phase 1: Setup and PostgreSQL User Model
- [x] Task: Create PostgreSQL persistence model `internal/infrastructure/persistence/postgres/user_model.go` 1ce4d7a
    - [x] Define `UserModel` struct with fields `ID`, `Document`, `IsActive`, and `CreatedAt`
    - [x] Add `db:"..."` tags using `snake_case` mapping
    - [x] Implement `ToDomain()` method to convert `UserModel` to `domain.User`
- [x] Task: Standardize PostgreSQL Repository Factory 1ce4d7a
    - [x] Ensure any repository factory for User is named `NewUserRepository` (not `NewGORMUserRepository`)
- [x] Task: Conductor - User Manual Verification 'Phase 1: Setup and PostgreSQL User Model' (Protocol in workflow.md) 1ce4d7a

## Phase 2: DynamoDB Token Model
- [x] Task: Create DynamoDB persistence model `internal/infrastructure/persistence/dynamo/token_model.go` 36a44e6
    - [x] Define `TokenModel` struct with fields `PK` (CPF), `Token`, and `ExpiresAt` (int64)
    - [x] Add `dynamodbav:"..."` tags using `snake_case` mapping
- [x] Task: Standardize DynamoDB Repository Factory 36a44e6
    - [x] Ensure any repository factory for Token is named `NewTokenRepository` (not `NewDynamoTokenRepository`)
    - [x] Use clean names for repository methods (e.g., `Save` instead of `SaveDynamoToken`)
- [x] Task: Conductor - User Manual Verification 'Phase 2: DynamoDB Token Model' (Protocol in workflow.md) 36a44e6

## Phase 3: Final Review and Standardization
- [~] Task: Audit for Clean Architecture Compliance
    - [ ] Verify no database tags exist in `internal/core/domain/`
    - [ ] Verify repository interfaces in `internal/core/ports/` remain technology-agnostic
    - [ ] Ensure all factories and methods follow the standardized nomenclature
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Final Review and Standardization' (Protocol in workflow.md)
