# Implementation Plan: User Authentication Lambda

## Phase 1: Cleanup and Structure Refactoring [checkpoint: dafa9e7]
- [x] Task: Remove existing `go-gin-lambda` directory and artifacts. (140f03e)
- [x] Task: Initialize new folder structure following Clean Architecture. (ed6615a)
    - [x] Create `cmd/auth/`
    - [x] Create `internal/core/domain/`, `internal/core/usecases/`, `internal/core/ports/`
    - [x] Create `internal/adapters/handlers/`, `internal/adapters/repositories/`
    - [x] Create `internal/infrastructure/`
- [x] Task: Initialize `go.mod` and install dependencies (`aws-sdk-go-v2`, `golang-jwt`). (4aabc9d)
- [x] Task: Conductor - User Manual Verification 'Phase 1: Cleanup and Structure Refactoring' (Protocol in workflow.md)

## Phase 2: Domain and Core Logic [checkpoint: e7f2492]
- [x] Task: Define User and Token domain entities in `internal/core/domain/`. (82e683c)
    - [x] Follow RDS User model from reference repository.
    - [x] Define Token entity with `token_id` and `user_id`.
- [x] Task: Define Repository interfaces in `internal/core/ports/`. (7688671)
    - [x] `UserRepository` for RDS.
    - [x] `TokenRepository` for DynamoDB.
- [x] Task: Implement Auth UseCase in `internal/core/usecases/`. (b9f556f)
    - [x] Logic for CPF regex validation.
    - [x] Logic for user status check.
    - [x] Logic for JWT generation (24h expiration, include jti).
- [x] Task: Conductor - User Manual Verification 'Phase 2: Domain and Core Logic' (Protocol in workflow.md)

## Phase 3: Adapters and Infrastructure
- [x] Task: Create RDS Repository integration file in `internal/adapters/repositories/rds_user_repository.go`. (48447ee)
    - [x] Logic not implemented, only file structure and model mapping.
- [x] Task: Create DynamoDB Repository integration file in `internal/adapters/repositories/dynamo_token_repository.go`. (9111890)
    - [x] Table: `user-authentication-token`, Fields: `token_id`, `user_id`.
    - [x] Logic not implemented, only file structure and model mapping.
- [ ] Task: Implement Lambda Handler in `internal/adapters/handlers/auth_handler.go`.
    - [ ] Use standard AWS Lambda handler (no Gin adapter).
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Adapters and Infrastructure' (Protocol in workflow.md)

## Phase 4: Entry Point and Deployment
- [ ] Task: Implement `main.go` in `cmd/auth/`.
    - [ ] Dependency injection for repositories and usecases.
    - [ ] Lambda start logic.
- [ ] Task: Update `deploy.sh` and `Makefile` for LocalStack deployment.
    - [ ] Configure `deploy.sh` to target LocalStack endpoints.
    - [ ] Configure `Makefile` for LocalStack-only deployment commands.
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Entry Point and Deployment' (Protocol in workflow.md)
