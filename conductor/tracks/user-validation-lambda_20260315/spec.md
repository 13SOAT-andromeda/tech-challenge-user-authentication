# Specification: User Validation Lambda

## Overview
Implement an AWS Lambda function in Go using the Gin framework to handle user authentication validation. The service will validate CPF structure via regex, check user status in an RDS database (`garagedb`), and if valid/active, generate a JWT token and store it in DynamoDB (`user-auth-tokens`).

## Functional Requirements
1. **CPF Regex Validation**:
   - Receive the CPF from the request header.
   - Validate the structure using: `^(\d{3}\.\d{3}\.\d{3}\-\d{2})?$`.
2. **RDS User Search**:
   - If the CPF is valid, query the `users` table in the `garagedb` RDS database.
   - Confirm if the user exists and is currently active.
3. **JWT Generation & DynamoDB Storage**:
   - If the user is found and active, generate a JWT token.
   - Store the generated token in the `user-auth-tokens` DynamoDB table.
4. **Response Handling**:
   - Return a successful response with the JWT token for active users.
   - Return appropriate error responses for invalid CPFs, non-existent users, or inactive users.

## Non-Functional Requirements
- **Simplicity & Consistency**: Follow the architecture of the `tech-challenge-s1` repository but eliminate unnecessary complexities.
- **Low Latency**: Ensure the validation and storage flow is fast enough for authentication use cases.
- **Reliability**: Handle database connections and AWS SDK calls with proper error management.

## Acceptance Criteria
- [ ] CPF structure is correctly validated against the specified regex.
- [ ] Only active users from the `users` table in `garagedb` are authenticated.
- [ ] JWT tokens are successfully generated and stored in the `user-auth-tokens` DynamoDB table.
- [ ] The API returns the correct HTTP status codes and JSON error messages for all failure scenarios.
