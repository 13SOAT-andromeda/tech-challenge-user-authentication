## Context

The application authenticates users and issues JWTs. The design calls for storing sessions in DynamoDB (table `user-auth-tokens`) using the JTI (a UUID v4) as the partition key, enabling server-side token revocation. Three bugs prevent this from working:

1. `SessionModel.UserID uint` is incorrectly tagged as `dynamodbav:"token_id"` — the partition key should be the JTI string.
2. `internal/service/session.go` is a stateless stub; `Create` never calls `PutItem` and `GetByID` returns `"not implemented"`.
3. `cmd/main.go` creates `SessionRepository` but discards it (`_ = ...`), wiring an unconnected stub to `AuthUseCase` instead.

The PostgreSQL `models` table is a remnant from a previous design and should be removed; sessions are DynamoDB-only.

## Goals / Non-Goals

**Goals:**
- Sessions are persisted to DynamoDB on every successful login
- Sessions are retrieved from DynamoDB on every token validation
- `SessionModel` schema matches the DynamoDB table key schema (`token_id` = JTI string)
- `SessionRepository` is properly injected into the session service and use case
- PostgreSQL no longer stores session data

**Non-Goals:**
- Session revocation endpoint (out of scope for this fix)
- TTL-based auto-expiry in DynamoDB (can be added separately)
- Refresh token rotation

## Decisions

### D1: Session service owns DynamoDB persistence
**Decision**: `internal/service/session.go` receives a `SessionRepository` dependency and calls it directly for both `Create` and `GetByID`.

**Why**: The service layer is the right owner of the persistence contract. The `SessionRepository` already implements the DynamoDB calls; we only need to inject it.

**Alternative considered**: Move persistence into the use case — rejected because it would bypass the service layer abstraction.

---

### D2: `SessionRepository` interface split: `Save` + `FindBySessionID`
**Decision**: Add `FindBySessionID(ctx, string) (*model.SessionModel, error)` to `SessionRepository`, backed by `GetItem`. Keep `Save` as-is.

**Why**: `GetByID` in the service needs a lookup by JTI string. Adding a typed method to the repository keeps the interface minimal and explicit.

**Alternative considered**: Expose a generic `Get` returning `map[string]types.AttributeValue` — rejected, too low-level.

---

### D3: Fix `SessionModel` field types
**Decision**: Replace `UserID uint` (tagged `token_id`) with `SessionID string` (tagged `token_id`) and add `UserID string` (tagged `user_id`).

**Why**: The DynamoDB table was created with `token_id` as a String type. Storing a `uint` there would cause a type mismatch. The JTI is a UUID string and must be the partition key.

---

### D4: Drop `models` table from PostgreSQL
**Decision**: Remove the `models` table from `db.AutoMigrate` and from the application entirely.

**Why**: Sessions are DynamoDB-only per the current design. Keeping a dead table causes confusion and wastes schema migration cycles.

## Risks / Trade-offs

- **Risk**: Existing test mocks use the old `SessionModel` shape → **Mitigation**: Update all test files that reference `SessionModel` or the stub service.
- **Risk**: DynamoDB connectivity issues during Lambda cold start (endpoint misconfiguration) → **Mitigation**: Already handled by `AWS_ENDPOINT_URL` env var in Lambda config.
- **Risk**: `FindBySessionID` returns `nil, nil` for missing sessions — caller must handle both error and nil session → **Mitigation**: `auth_usecase.go` already checks `session == nil` before proceeding.

## Migration Plan

1. Apply code changes (model, repository, service, main.go)
2. Rebuild Lambda binary (`make build`)
3. Redeploy to LocalStack (`make deploy` or equivalent)
4. Re-run seed to ensure test user exists with known password hash
5. Smoke test: `POST /sessions` → verify DynamoDB table contains session item
6. No rollback concerns — DynamoDB table schema is backward compatible (String partition key unchanged)
