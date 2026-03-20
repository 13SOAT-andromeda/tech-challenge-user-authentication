## Why

User sessions are silently lost on every login because the `SessionRepository` (wired to DynamoDB) is created but immediately discarded in `main.go`, and the session service is an in-memory stub that never persists or retrieves data. This means the token validation flow (`GetByID`) always returns `"not implemented"`, breaking the JWT revocation mechanism entirely.

## What Changes

- **Fix `SessionModel`**: rename `UserID uint` (mapped to `token_id`) to `SessionID string` so the DynamoDB partition key correctly stores the JTI UUID.
- **Fix `SessionRepository`**: add `GetItem` support (`FindBySessionID`) so sessions can be retrieved by JTI for token validation.
- **Fix `internal/service/session.go`**: replace the in-memory stub with a real implementation that injects and uses the `SessionRepository` for both `Create` (PutItem) and `GetByID` (GetItem).
- **Fix `cmd/main.go`**: wire `SessionRepository` into `sessionSvc` instead of discarding it.
- **Remove PostgreSQL session table**: the `models` table (a legacy session store) should be dropped — sessions are owned by DynamoDB only.
- **Update seed**: document the correct test credentials so the seed password hash is reproducible and testable.

## Capabilities

### New Capabilities

- `session-persistence`: Persist and retrieve user sessions in DynamoDB using the JTI as the partition key, enabling JWT revocation via session lookup.

### Modified Capabilities

<!-- None — no existing spec files exist; this is net-new spec coverage. -->

## Impact

- **`internal/adapters/database/model/session.go`** — model field rename
- **`internal/adapters/database/repositories/session_repository.go`** — add `GetItem` interface method and `FindBySessionID` method
- **`internal/service/session.go`** — inject repository, implement `Create` + `GetByID`
- **`cmd/main.go`** — wire `sessionRepo` → `sessionSvc`
- **`scripts/seed.sql`** — regenerate password hash with a known value (`admin123`)
- **PostgreSQL `models` table** — drop or ignore (sessions move to DynamoDB)
- No API contract changes; no external dependencies added
