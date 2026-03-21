# Product Guidelines

## Prose Style
- **Concise & Technical**: Focus on technical accuracy and directness, ideal for API documentation and internal developer communication.

## Branding and Naming
- **Idiomatic Go Style**: Adhere to standard Go conventions, using `MixedCaps` for exported identifiers and `camelCase` for internal names.

## UX Principles
- **Clear Error Messages**: Provide detailed and actionable error messages in all API responses to simplify debugging for consumers.
- **Payload Efficiency**: Prioritize minimal response payload size for fast parsing and efficient network usage.
- **Self-Documenting API**: Ensure all public endpoints are documented with clear examples and intuitive naming.

## Error Handling
- **Detailed JSON Errors**: Include custom internal error codes and descriptive text in the JSON body for all failure scenarios.
