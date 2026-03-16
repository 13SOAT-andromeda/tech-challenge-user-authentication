# Specification: User Authentication Lambda

## Overview
Implement an AWS Lambda function in Go using the Gin framework to handle user authentication validation. The service will validate the CPF format, verify the user's status in an RDS database, and generate a JWT token for active users, storing it in DynamoDB.

## Functional Requirements
1.  **CPF Format Validation**: Validate the CPF received in the header using the regex `^(\d{3}\.\d{3}\.\d{3}\-\d{2})?$`.
2.  **RDS User Query**: Search for the user in the `users` table of the `garagedb` RDS database.
3.  **User Status Validation**: Verify if the user is active (true/false).
4.  **JWT Token Generation**: Generate a JWT token with a 24-hour expiration for active users.
5.  **Token Persistence**: Store the generated JWT token in the `user-auth-tokens` DynamoDB table.

## Acceptance Criteria
-   Successful validation of a valid and active CPF returns a 24-hour JWT token.
-   Invalid CPF format returns a 400 Bad Request.
-   Inactive or non-existent users return a 401 Unauthorized.
-   Tokens are correctly persisted in DynamoDB for subsequent verification.
