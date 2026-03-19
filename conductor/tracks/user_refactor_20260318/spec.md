# Specification: Refactor User Model and Repository

## Overview
Refactor the `tech-challenge-user-validation` Lambda to align its `User` data structure and repository implementation with the provided reference architecture (`database/model/user/user.go` and `database/repository/user_repository.go`). This includes adopting generic naming conventions, integrating the embedded `Address` model, and updating tests.

## Functional Requirements
- **Model Synchronization:**
  - Implement the `user.Model` with fields: `Name`, `Email`, `Contact`, `Address`, `Password`, `Role`.
  - Implement the `address.Model` to support the embedded address logic in the `User` struct.
  - Implement `ToDomain` and `FromDomain` mapping methods for both models, utilizing the provided `encryption` and `domain` packages.
- **Repository Implementation:**
  - Create a generic `BaseRepository`.
  - Implement `userRepository` incorporating `Search` and `GetByEmail` methods.
- **Initialization:**
  - Update the Lambda initialization to run `gorm.AutoMigrate` for the new `user.Model` (and `address.Model`).
- **Naming Standardization:**
  - Remove specific implementation names (e.g., `NewGORMUserRepository`).
  - Align interfaces with standard `ports.UserRepository` conventions.

## Non-Functional Requirements
- **Code Quality:** Ensure idiomatic Go code, applying Clean Code principles and maintaining the current directory structure.
- **Testability:** The use cases should be decoupled from the repository logic.

## Acceptance Criteria
- [ ] `user.Model` and `address.Model` are implemented and match the reference structure.
- [ ] `userRepository` is implemented using `BaseRepository` and successfully executes `GetByEmail` and `Search`.
- [ ] Lambda initialization runs AutoMigrate for the new models.
- [ ] Unit tests are written using explicit mock structs for repositories, achieving >80% coverage on validation/login logic.
- [ ] No specific implementation names leak into interface or constructor names.

## Out of Scope
- Modifying the underlying RDS instance configuration beyond `AutoMigrate`.
- Changes to the DynamoDB token/session implementation unless strictly necessary for compatibility.