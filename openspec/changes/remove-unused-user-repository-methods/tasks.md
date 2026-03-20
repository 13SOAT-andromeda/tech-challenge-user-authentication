## 1. Interface and Struct Cleanup

- [x] 1.1 Remove `GetByEmail` method from `UserRepository` interface in `internal/core/ports/user_repository.go`
- [x] 1.2 Remove `Search` method from `UserRepository` interface in `internal/core/ports/user_repository.go`
- [x] 1.3 Remove `UserSearch` struct from `internal/core/ports/user_repository.go`

## 2. Implementation Cleanup

- [x] 2.1 Remove `GetByEmail` method implementation from `internal/adapters/database/repositories/user_repository.go`
- [x] 2.2 Remove `Search` method implementation from `internal/adapters/database/repositories/user_repository.go`

## 3. Test Mock Cleanup

- [x] 3.1 Remove `GetByEmail` from `mockUserRepository` in `internal/core/usecases/auth_usecase_test.go`
- [x] 3.2 Remove `Search` from `mockUserRepository` in `internal/core/usecases/auth_usecase_test.go`
- [x] 3.3 Remove `GetByEmail` from `mockUserRepository` in `internal/adapters/http/handlers/auth_handler_test.go`
- [x] 3.4 Remove `Search` from `mockUserRepository` in `internal/adapters/http/handlers/auth_handler_test.go`

## 4. Repository Test Cleanup

- [x] 4.1 Remove `GetByEmail` test cases from `internal/adapters/database/repositories/user_repository_test.go`
- [x] 4.2 Remove `Search` test cases from `internal/adapters/database/repositories/user_repository_test.go`

## 5. Verification

- [x] 5.1 Run `go build ./...` and confirm no compilation errors
- [x] 5.2 Run `go test ./...` and confirm all tests pass
