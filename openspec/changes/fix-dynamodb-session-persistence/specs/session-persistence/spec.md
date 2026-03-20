## ADDED Requirements

### Requirement: Session is persisted to DynamoDB on login
On a successful login, the system SHALL create a session record in the DynamoDB `user-auth-tokens` table using the JTI (UUID v4) as the partition key (`token_id`), storing the user ID and expiration timestamp.

#### Scenario: Successful login creates DynamoDB session
- **WHEN** a user submits valid credentials (document + password)
- **THEN** a session item SHALL exist in DynamoDB with `token_id` equal to the JTI from the returned JWT

#### Scenario: Session stores correct user and expiration
- **WHEN** a session is created
- **THEN** the DynamoDB item SHALL contain `user_id` matching the authenticated user's ID and `expires_at` set to 7 days from login time

#### Scenario: Failed login does not create session
- **WHEN** a user submits invalid credentials
- **THEN** no session item SHALL be written to DynamoDB

---

### Requirement: Session is retrievable by JTI for token validation
The system SHALL retrieve a session from DynamoDB using the JTI claim from the JWT during token validation.

#### Scenario: Valid token with existing session passes validation
- **WHEN** a token validation request is made with a JWT whose JTI exists in DynamoDB
- **THEN** the system SHALL return validation success

#### Scenario: Token with missing session fails validation
- **WHEN** a token validation request is made with a JWT whose JTI does NOT exist in DynamoDB
- **THEN** the system SHALL return `"session not found or revoked"`

---

### Requirement: SessionModel partition key is the JTI string
The `SessionModel` struct SHALL map the `token_id` DynamoDB attribute to a `string` field representing the JTI (UUID v4), not a `uint` user ID.

#### Scenario: Session item written with string token_id
- **WHEN** a session is saved to DynamoDB
- **THEN** the `token_id` attribute SHALL be a String type containing the UUID JTI

---

### Requirement: Sessions are not stored in PostgreSQL
The system SHALL NOT store session data in PostgreSQL. The `models` table SHALL be removed from the auto-migration list.

#### Scenario: Login does not write to PostgreSQL sessions table
- **WHEN** a user logs in successfully
- **THEN** no session row SHALL be inserted into PostgreSQL
