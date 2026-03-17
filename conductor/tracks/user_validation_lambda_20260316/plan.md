# Implementation Plan: User Validation Lambda

## Phase 1: Domain & Ports Audit (Already partially implemented)
- [ ] Task: Define Domain Models (User, Token)
    - [x] Create `internal/core/domain/user.go`
    - [x] Create `internal/core/domain/token.go`
- [ ] Task: Define Repository Ports
    - [x] Create `internal/core/ports/user_repository.go` (RDS)
    - [x] Create `internal/core/ports/token_repository.go` (DynamoDB)
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Domain & Ports Audit' (Protocol in workflow.md)

## Phase 2: UseCase & Domain Enhancement (TDD)
- [ ] Task: Update `AuthUseCase` and `domain.User` for precision
    - [ ] Ensure `domain.User` fields match the specification (e.g., `Document` vs `CPF`)
    - [ ] Update `Authenticate` method to return expiration time
    - [ ] Write failing tests for the new response format
- [ ] Task: Update `AuthHandler` for spec compliance
    - [ ] Update header extraction to use `x-user-cpf`
    - [ ] Update response marshalling to include `expires_at`
- [ ] Task: Conductor - User Manual Verification 'Phase 2: UseCase & Domain Enhancement' (Protocol in workflow.md)

## Phase 3: Infrastructure - PostgreSQL RDS Repository (TDD)
- [ ] Task: Implement RDS Repository with PostgreSQL Driver
    - [ ] Add `github.com/lib/pq` and configure connection pool
    - [ ] Write failing tests for `GetByDocument`
    - [ ] Implement secure logic with Prepared Statements
    - [ ] Manage credentials via environment variables (`DB_HOST`, `DB_USER`, `DB_PASSWORD`)
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Infrastructure - PostgreSQL RDS Repository' (Protocol in workflow.md)

## Phase 4: Infrastructure - DynamoDB Token Repository (TDD)
- [ ] Task: Implement DynamoDB Repository with AWS SDK v2
    - [ ] Setup AWS SDK v2 configuration
    - [ ] Write failing tests for `Save`
    - [ ] Implement `Save` logic targeting `user-auth-tokens` table (Partition Key: `token_hash`)
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Infrastructure - DynamoDB Token Repository' (Protocol in workflow.md)

## Phase 5: Final Integration & Deployment (LocalStack)
- [ ] Task: Complete `main.go` and Deployment Scripts
    - [ ] Setup all environment variables and dependency injection
    - [ ] Create `init-aws.sh` and `init-db.sql` for LocalStack setup
    - [ ] Update `deploy.sh` and `Makefile`
- [ ] Task: Deployment and Manual Verification
    - [ ] Build and deploy to LocalStack
    - [ ] Verify full flow with `curl`
- [ ] Task: Conductor - User Manual Verification 'Phase 5: Final Integration & Deployment' (Protocol in workflow.md)
