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
- [ ] Task: Create DynamoDB persistence model `internal/infrastructure/persistence/dynamo/token_model.go`
    - [ ] Define `TokenModel` struct with fields `PK` (CPF), `Token`, and `ExpiresAt` (int64)
    - [ ] Add `dynamodbav:"..."` tags using `snake_case` mapping
- [ ] Task: Standardize DynamoDB Repository Factory
    - [ ] Ensure any repository factory for Token is named `NewTokenRepository` (not `NewDynamoTokenRepository`)
    - [ ] Use clean names for repository methods (e.g., `Save` instead of `SaveDynamoToken`)
- [ ] Task: Conductor - User Manual Verification 'Phase 2: DynamoDB Token Model' (Protocol in workflow.md)

## Phase 3: Final Review and Standardization
- [ ] Task: Audit for Clean Architecture Compliance
    - [ ] Verify no database tags exist in `internal/core/domain/`
    - [ ] Verify repository interfaces in `internal/core/ports/` remain technology-agnostic
    - [ ] Ensure all factories and methods follow the standardized nomenclature
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Final Review and Standardization' (Protocol in workflow.md)
