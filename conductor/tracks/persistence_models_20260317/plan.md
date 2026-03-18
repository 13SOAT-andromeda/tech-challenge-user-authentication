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
- [x] Task: Audit for Clean Architecture Compliance 78e10c3
    - [x] Verify no database tags exist in `internal/core/domain/`
    - [x] Verify repository interfaces in `internal/core/ports/` remain technology-agnostic
    - [x] Ensure all factories and methods follow the standardized nomenclature
- [x] Task: Conductor - User Manual Verification 'Phase 3: Final Review and Standardization' (Protocol in workflow.md) 78e10c3

## Phase 4: Refactor Directory Structure and Filenames (User Suggestion)
- [x] Task: Move and rename persistence models
    - [x] Move `internal/infrastructure/persistence/postgres/user_model.go` to `internal/adapters/database/model/user.go` (and update tests)
    - [x] Move `internal/infrastructure/persistence/dynamo/token_model.go` to `internal/adapters/database/model/token.go` (and update tests)
- [x] Task: Move repository implementations
    - [x] Move `internal/adapters/repositories/user_repository.go` to `internal/adapters/database/repositories/user_repository.go` (and update tests)
    - [x] Move `internal/adapters/repositories/token_repository.go` to `internal/adapters/database/repositories/token_repository.go` (and update tests)
- [x] Task: Update all package imports and references
    - [x] Update imports across the project to match the new `model` and `repositories` paths.
- [~] Task: Verify all tests pass
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Refactor Directory Structure and Filenames' (Protocol in workflow.md)
