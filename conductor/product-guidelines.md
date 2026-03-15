# Product Guidelines

## Prose Style
- **Concise & Technical**: All documentation, including API responses and logs, should prioritize technical accuracy and directness. Avoid unnecessary filler or overly narrative explanations.

## Branding and Structure
- **Consistent with `tech-challenge-s1`**: Adhere to the naming conventions, project organization, and structural patterns established in the `https://github.com/13SOAT-andromeda/tech-challenge-s1` repository.
- **Minimalist Implementation**: While following the `s1` structure, aim for the simplest and most consistent implementation possible, eliminating unnecessary complexities.

## UX Principles
- **API-First Design**: Prioritize clear, actionable API responses. Every endpoint should be easy to understand and integrate with.
- **Performance-Optimized**: Ensure low-latency validation and token generation to minimize impact on downstream services.
- **Simplicity and Ease of Use**: Maintain a clean API interface with minimal required inputs to simplify integration for internal developers.

## Error Handling
- **Detailed Custom Errors**: Provide clear, custom error codes and descriptive messages in the JSON response body. This helps developers quickly diagnose issues while maintaining a consistent error format.
