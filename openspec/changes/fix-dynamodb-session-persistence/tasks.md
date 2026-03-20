## 1. Fix SessionModel Schema

- [x] 1.1 In `internal/adapters/database/model/session.go`, replace `UserID uint` (tagged `dynamodbav:"token_id"`) with `SessionID string` (tagged `dynamodbav:"token_id"`) and add `UserID string` (tagged `dynamodbav:"user_id"`)
- [x] 1.2 Remove the unused `RefreshToken *string` field from `SessionModel` (sessions do not need to store the refresh token value)

## 2. Fix SessionRepository — Add GetItem Support

- [x] 2.1 Add `GetItem` to the `SessionDynamoClient` interface in `internal/adapters/database/repositories/session_repository.go`
- [x] 2.2 Add `FindBySessionID(ctx context.Context, sessionID string) (*model.SessionModel, error)` method to `SessionRepository` using `GetItem` with the `token_id` key
- [x] 2.3 Update any existing tests in `session_repository_test.go` to reflect the new model shape and added method

## 3. Fix Session Service — Wire Repository and Implement Persistence

- [x] 3.1 Define a `sessionRepository` interface in `internal/service/session.go` with `Save` and `FindBySessionID` methods matching the `SessionRepository`
- [x] 3.2 Update `sessionService` struct to hold a `repo sessionRepository` dependency
- [x] 3.3 Update `NewSessionService` to accept the repository and store it
- [x] 3.4 Implement `Create`: validate inputs, build `model.SessionModel`, call `repo.Save`, and return `*ports.Session`
- [x] 3.5 Implement `GetByID`: call `repo.FindBySessionID`, return `nil, nil` if not found, map model to `*ports.Session` if found

## 4. Wire Everything in main.go

- [x] 4.1 Remove `_ =` from `repositories.NewSessionRepository(...)` and assign it to `sessionRepo`
- [x] 4.2 Pass `sessionRepo` to `services.NewSessionService(sessionRepo)`
- [x] 4.3 Remove `models` struct (or its import) from `db.AutoMigrate` call if it references a PostgreSQL session table

## 5. Update Seed and Test Credentials

- [x] 5.1 Regenerate the bcrypt hash in `scripts/seed.sql` for password `admin123` using Go's `bcrypt.DefaultCost` (`$2a$10$...`) so it is testable
- [x] 5.2 Verify the seed user can log in via the Lambda invoke test

## 6. Rebuild and Smoke Test

- [x] 6.1 Run `go build ./...` (or `make build`) and confirm zero compilation errors
- [x] 6.2 Run unit tests: `go test ./...`
- [x] 6.3 Redeploy Lambda to LocalStack (`make deploy` or equivalent)
- [x] 6.4 Invoke login via Lambda or API Gateway and confirm HTTP 200 with JWT tokens
- [x] 6.5 Run `awslocal dynamodb scan --table-name user-auth-tokens` and confirm a session item exists with the correct `token_id` (UUID), `user_id`, and `expires_at`
