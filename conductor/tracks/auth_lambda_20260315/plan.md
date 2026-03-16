# Implementation Plan: User Authentication Lambda

## Phase 1: Cleanup and Structure Refactoring
- [x] Task: Remove existing `go-gin-lambda` directory and artifacts. (140f03e)
- [ ] Task: Initialize new folder structure following Clean Architecture.
    - [ ] Create `cmd/auth/`
    - [ ] Create `internal/core/domain/`, `internal/core/usecases/`, `internal/core/ports/`
    - [ ] Create `internal/adapters/handlers/`, `internal/adapters/repositories/`
    - [ ] Create `internal/infrastructure/`
- [ ] Task: Initialize `go.mod` and install dependencies (`aws-sdk-go-v2`, `golang-jwt`).
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Cleanup and Structure Refactoring' (Protocol in workflow.md)

## Phase 2: Domain and Core Logic
- [ ] Task: Define User and Token domain entities in `internal/core/domain/`.
    - [ ] Follow RDS User model from reference repository.
    - [ ] Define Token entity with `token_id` and `user_id`.
- [ ] Task: Define Repository interfaces in `internal/core/ports/`.
    - [ ] `UserRepository` for RDS.
    - [ ] `TokenRepository` for DynamoDB.
- [ ] Task: Implement Auth UseCase in `internal/core/usecases/`.
    - [ ] Logic for CPF regex validation.
    - [ ] Logic for user status check.
    - [ ] Logic for JWT generation (24h expiration, include jti).
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Domain and Core Logic' (Protocol in workflow.md)

## Phase 3: Adapters and Infrastructure
- [ ] Task: Create RDS Repository integration file in `internal/adapters/repositories/rds_user_repository.go`.
    - [ ] Logic not implemented, only file structure and model mapping.
- [ ] Task: Create DynamoDB Repository integration file in `internal/adapters/repositories/dynamo_token_repository.go`.
    - [ ] Table: `user-authentication-token`, Fields: `token_id`, `user_id`.
    - [ ] Logic not implemented, only file structure and model mapping.
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
