# Specification: Session Management and JTI Integration (Minimalist + user_id)

## Overview
Refactor the session management logic within the `tech-challenge-user-validation` service to incorporate JTI (JSON Token Identifier) and persist session metadata in DynamoDB. This ensures that JWTs can be tracked and validated as "active" in real-time, enabling immediate revocation and improving overall security.

## Functional Requirements
- **Login Flow:**
    - Generate a unique JTI (UUID v4) for both Access and Refresh Tokens.
    - Persist the JTI, `user_id`, and `expires_at` in DynamoDB using the existing `sessionService`.
    - Include the `jti` claim in the generated JWTs.
- **Validation Flow:**
    - The `Validate` method in the Use Case must extract the `jti` from the token's claims.
    - Check the status of the session in DynamoDB via `sessionService`.
    - Mark the validation as successful only if the session is found in the database.
- **Code Standards:**
    - Rename implementation-specific names to generic port-aligned names (e.g., `UserRepository`, `TokenRepository`).
    - Maintain Clean Architecture principles and dependency injection patterns.

## Non-Functional Requirements
- **Security:** Ensure sessions in DynamoDB follow the same TTL as the Refresh Token.
- **Performance:** DynamoDB queries for JTI validation should be efficient to minimize latency in the authentication flow.

## Acceptance Criteria
- [ ] Login returns Access and Refresh Tokens containing a UUID JTI.
- [ ] A record is created in DynamoDB for every new session containing only `pk` (JTI), `user_id`, and `expires_at`.
- [ ] The `Validate` method correctly rejects tokens if the corresponding JTI is missing in DynamoDB.
- [ ] Repository and service names are refactored to be generic and consistent.

## Out of Scope
- Creation of a separate validation Lambda entry point (`cmd/validate/main.go`). This track focuses on the core Use Case logic.
- Implementation of the `ports.SessionService` (assumed to be already existing for integration).