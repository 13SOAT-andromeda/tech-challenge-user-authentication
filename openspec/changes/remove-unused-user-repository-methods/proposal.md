## Why

The `UserRepository` interface defines `GetByEmail` and `Search` methods that have concrete implementations and tests but are never called anywhere in the production codebase. These dead code paths increase maintenance burden, inflate the interface contract, and require mock updates in every test that uses `UserRepository`.

## What Changes

- Remove `GetByEmail` method from the `UserRepository` interface
- Remove `GetByEmail` implementation from the concrete `userRepository` struct
- Remove `Search` method from the `UserRepository` interface
- Remove `Search` implementation from the concrete `userRepository` struct
- Remove `UserSearch` struct (used only by `Search`)
- Remove `GetByEmail` and `Search` from all test mock implementations of `UserRepository`
- Remove tests for the deleted methods in `user_repository_test.go`

## Capabilities

### New Capabilities
<!-- None — this is a cleanup change with no new behaviors -->

### Modified Capabilities
- `user-repository`: Interface contract is reduced; `GetByEmail` and `Search` are no longer part of the public API

## Impact

- **Interface**: `internal/core/ports/user_repository.go` — interface shrinks from 3 to 1 method; `UserSearch` struct removed
- **Implementation**: `internal/adapters/database/repositories/user_repository.go` — two methods removed
- **Tests**: `internal/adapters/database/repositories/user_repository_test.go` — test cases for removed methods deleted
- **Mocks**: `internal/core/usecases/auth_usecase_test.go` and `internal/adapters/http/handlers/auth_handler_test.go` — mock structs trimmed
- No API surface change; no external dependencies affected
