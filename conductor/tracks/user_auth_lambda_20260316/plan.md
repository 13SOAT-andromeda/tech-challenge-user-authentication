# Implementation Plan: User Auth Lambda (Go/GORM/DynamoDB)

## Phase 1: Setup and Infrastructure
- [x] Task: Define Domain Models (User, Token) ffffef9
    - [x] Create `internal/core/domain/user.go`
    - [x] Create `internal/core/domain/token.go`
- [ ] Task: Create LocalStack Initialization Scripts
    - [ ] Create `scripts/init-rds.sql` to setup `garagedb` and users table
    - [ ] Create `scripts/init-dynamodb.sh` to setup `user-auth-tokens` table
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Setup and Infrastructure' (Protocol in workflow.md)

## Phase 2: Repository Layer (TDD)
- [ ] Task: Implement PostgreSQL Repository (GORM)
    - [ ] Write failing tests for `GetByDocument` in `internal/adapters/repositories/user_repository_test.go`
    - [ ] Implement `GetByDocument` in `internal/adapters/repositories/user_repository.go` using `gorm.io/gorm` and `gorm.io/driver/postgres`
- [ ] Task: Implement DynamoDB Repository
    - [ ] Write failing tests for `SaveToken` in `internal/adapters/repositories/token_repository_test.go`
    - [ ] Implement `SaveToken` in `internal/adapters/repositories/token_repository.go` using AWS SDK v2
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Repository Layer (TDD)' (Protocol in workflow.md)

## Phase 3: Usecase Layer (TDD)
- [ ] Task: Implement Authentication Usecase
    - [ ] Write failing tests for `Authenticate` in `internal/core/usecases/auth_usecase_test.go`
    - [ ] Implement `Authenticate` in `internal/core/usecases/auth_usecase.go` (Regex validation, GORM lookup, JWT generation, DynamoDB persistence)
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Usecase Layer (TDD)' (Protocol in workflow.md)

## Phase 4: Handler Layer (TDD)
- [ ] Task: Implement Lambda Handler
    - [ ] Write failing tests for header extraction and status codes in `internal/adapters/handlers/auth_handler_test.go`
    - [ ] Implement the handler in `internal/adapters/handlers/auth_handler.go`
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Handler Layer (TDD)' (Protocol in workflow.md)

## Phase 5: Integration and Verification
- [ ] Task: Manual Verification
    - [ ] Test success flow (200 OK with token)
    - [ ] Test invalid CPF format (400 Bad Request)
    - [ ] Test user not found/inactive (404 Not Found)
- [ ] Task: Conductor - User Manual Verification 'Phase 5: Integration and Verification' (Protocol in workflow.md)
