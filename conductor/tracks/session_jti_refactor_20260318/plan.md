# Implementation Plan - Session Management and JTI Integration (Minimalist + user_id)

Refactor session management to use UUID-based JTI and persist session metadata in DynamoDB for real-time validation.

## Phase 1: Preparation and Generic Refactoring
Goal: Align naming conventions with generic ports and ensure the environment is ready.

- [x] Task: Refactor `UserRepository` and `TokenRepository` to remove implementation-specific prefixes (e.g., `GORM`, `Dynamo`). 26c005b
    - [x] Rename files and structs in `internal/adapters/database/repositories/`. 26c005b
    - [x] Update constructor names (e.g., `NewUserRepository`). 26c005b
    - [x] Update all references in `cmd/auth/main.go`. 26c005b
- [x] Task: Ensure `ports.SessionService` is correctly defined and accessible. 26c005b
    - [x] Verify `internal/core/ports/session_service.go`. 26c005b
- [x] Task: Conductor - User Manual Verification 'Phase 1: Preparation and Generic Refactoring' (Protocol in workflow.md) 26c005b

## Phase 2: Login Flow (Minimalist Persistence)
Goal: Update the Login process to generate a UUID JTI and save minimalist session metadata.

- [x] Task: Update `AuthUseCase.Authenticate` (or `Login`) to generate a UUID v4 as JTI. 9479164
    - [x] Write failing test in `internal/core/usecases/auth_usecase_test.go` verifying JTI presence in claims. 9479164
    - [x] Implement UUID generation using a standard library (e.g., `google/uuid`). 9479164
    - [x] Include `jti` claim in both Access and Refresh Tokens. 9479164
- [x] Task: Integrate `sessionService.Create` in the Login flow. 9479164
    - [x] Write failing test verifying `sessionService.Create` is called with `pk` (JTI), `user_id`, and `expires_at`. 9479164
    - [x] Implement the call to `sessionService.Create` before returning the tokens. 9479164
    - [x] Ensure `expires_at` matches the Refresh Token expiry. 9479164
- [x] Task: Conductor - User Manual Verification 'Phase 2: Login Flow (Minimalist Persistence)' (Protocol in workflow.md) 9479164

## Phase 3: Validation Flow (JTI Verification)
Goal: Implement real-time session validation using the JTI from the token.

- [x] Task: Update `AuthUseCase.Validate` to extract and verify JTI. 9479164
    - [x] Write failing test for `Validate` where a token with an invalid/revoked JTI is rejected. 9479164
    - [x] Implement JTI extraction from JWT claims. 9479164
    - [x] Call `sessionService.GetByID` (using JTI) to verify session exists in DynamoDB. 9479164
    - [x] Ensure the validation fails if the session is not found. 9479164
- [x] Task: Conductor - User Manual Verification 'Phase 3: Validation Flow (JTI Verification)' (Protocol in workflow.md) 9479164

## Phase 4: Final Integration and Cleanup
Goal: Ensure everything is tied together and following the style guides.

- [ ] Task: Perform final end-to-end manual test in LocalStack.
    - [ ] Verify DynamoDB records (`pk`, `user_id`, `expires_at`) after a successful login.
    - [ ] Verify `Validate` behavior with active and "revoked" sessions.
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Final Integration and Cleanup' (Protocol in workflow.md)