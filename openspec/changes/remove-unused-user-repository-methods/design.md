## Context

The `UserRepository` interface currently exposes three methods: `GetByDocument`, `GetByEmail`, and `Search`. Only `GetByDocument` is called in production code (the `Login` use case). `GetByEmail` and `Search` — along with the `UserSearch` helper struct — have full implementations and tests but are dead code. This cleanup reduces the interface surface to only what is actually needed.

## Goals / Non-Goals

**Goals:**
- Delete `GetByEmail` and `Search` from the interface, implementation, and all test mocks
- Remove the now-orphaned `UserSearch` struct
- Delete the corresponding test cases from `user_repository_test.go`
- Ensure the remaining `GetByDocument` path is unaffected and tests stay green

**Non-Goals:**
- Introducing any new repository methods
- Changing the behavior of `GetByDocument` or the `Login` flow
- Modifying the database schema or DynamoDB table configuration

## Decisions

**Remove both methods entirely (vs. deprecating with a comment)**
Keeping dead code with a comment still adds noise to mocks and tests. Since this is an internal port (not a public SDK), there are no external consumers to protect. Hard removal is the right call.

**Delete the tests for removed methods (vs. leaving them)**
Tests for deleted code are misleading — they pass trivially if the code is gone, or they break. Removing them keeps the test suite honest.

**No interface versioning or adapter needed**
All consumers of `UserRepository` are internal. A simple find-and-delete is sufficient; no compatibility shim is required.

## Risks / Trade-offs

- [Someone planned to use GetByEmail/Search in a near-future feature] → Check with teammates before merging; if needed, the methods can be re-added from git history
- [Tests were covering behavior that will be wanted again] → Git history preserves the implementation; restoration is trivial

## Migration Plan

1. Delete methods from interface (`ports/user_repository.go`)
2. Delete implementations from concrete struct (`repositories/user_repository.go`)
3. Update mocks in `auth_usecase_test.go` and `auth_handler_test.go`
4. Delete test cases in `user_repository_test.go`
5. Run `go build ./...` and `go test ./...` to confirm zero regressions
6. Rollback: `git revert` the commit if a consumer appears post-merge
