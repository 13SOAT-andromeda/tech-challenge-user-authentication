# Tech Stack

## Language
- **Go v1.25.1**: The core programming language for the application logic.

## Frameworks
- **Gin v1.11.0**: Used for HTTP routing and API management within the Lambda environment.

## Infrastructure & Cloud
- **AWS Lambda**: The serverless compute platform for the application.
- **AWS RDS (garagedb)**: Used for persistent storage of user data.
- **Amazon DynamoDB (user-auth-tokens)**: Used for high-speed, scalable session and token storage.
- **LocalStack**: For local emulation of AWS services during development and testing.

## Deployment & Tooling
- **Docker/Docker Compose**: Manages the LocalStack environment.
- **deploy.sh**: Custom shell script for building, zipping, and deploying the Lambda function.
