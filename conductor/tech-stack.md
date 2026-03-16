# Tech Stack

## Language
- **Go (v1.23.0)**: The core programming language for the API, chosen for its performance and native support in AWS Lambda.

## Frameworks
- **Standard AWS Lambda Library**: Used for handling events directly without additional web framework overhead, ensuring minimal cold starts and simplicity.

## Infrastructure
- **LocalStack**: Provides a local emulation of AWS services for consistent development and testing.
- **Docker & Docker Compose**: Manages the LocalStack container and environment configuration.

## Deployment & Tooling
- **deploy.sh**: Custom shell script for building the Go binary and deploying the zipped function to LocalStack.
- **Makefile**: Provides common commands for local development, building, and deployment.
