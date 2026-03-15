# Initial Concept
A Go Gin Lambda designed to validate user CPF against RDS (GarageDB) and store JWT tokens in DynamoDB for authentication.

# Product Definition

## Vision
To provide a lightweight and reliable user validation service that checks CPF formatting, verifies user status in RDS, and generates authentication tokens stored in DynamoDB.

## Target Audience
- **Downstream Services**: Other microservices requiring reliable user validation.
- **End-Users (via API)**: Users attempting to authenticate and access protected resources.
- **Internal Developers**: Engineers looking for a standardized Go Gin Lambda architecture.

## Primary Goals
- **User Validation**: Ensure CPF is correctly formatted and exists as an active user in the `garagedb` RDS.
- **Authentication Flow**: Seamlessly transition from validation to JWT generation and session storage in DynamoDB.
- **High Availability**: Leverage AWS Lambda and DynamoDB for low-latency, scalable operations.

## Core Features
- **CPF Regex Validation**: Strict validation of the CPF structure using: `^(\d{3}\.\d{3}\.\d{3}\-\d{2})?$`
- **RDS User Check**: Query the `users` table in the `garagedb` database hosted on AWS RDS.
- **DynamoDB Token Storage**: Store generated JWT tokens in the `user-auth-tokens` table for session management.

## Success Criteria
- Valid CPFs of active users correctly result in a JWT token stored in DynamoDB.
- Invalid or inactive CPFs are rejected with appropriate error responses.
- The system handles database connection pooling and token lifecycle management efficiently.
