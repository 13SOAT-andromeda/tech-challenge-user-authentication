## REMOVED Requirements

### Requirement: GetByEmail
The system SHALL provide a method to retrieve a user by their email address.
**Reason**: `GetByEmail` is not called anywhere in production code. Retaining it inflates the interface contract and requires mock boilerplate in every test that uses `UserRepository`.
**Migration**: No migration needed; there are no callers. If a future use case requires lookup by email, re-add the method at that time.

#### Scenario: Retrieve existing user by email
- **WHEN** `GetByEmail` is called with a valid email
- **THEN** the corresponding user record is returned

#### Scenario: User not found by email
- **WHEN** `GetByEmail` is called with an email that does not exist
- **THEN** nil is returned with no error

### Requirement: Search users by name
The system SHALL provide a method to search users by name using a partial, case-insensitive match via `UserSearch` parameters.
**Reason**: `Search` is not called anywhere in production code. The `UserSearch` struct exists solely to support this method.
**Migration**: No migration needed; there are no callers. Re-add if a search use case is introduced.

#### Scenario: Search returns matching users
- **WHEN** `Search` is called with a non-empty name filter
- **THEN** all users whose name contains the filter string (case-insensitive) are returned

#### Scenario: Search returns empty list when no match
- **WHEN** `Search` is called with a name that matches no user
- **THEN** an empty slice is returned
