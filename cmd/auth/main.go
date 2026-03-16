package main

import (
	"os"

	"tech-challenge-user-validation/internal/adapters/handlers"
	"tech-challenge-user-validation/internal/adapters/repositories"
	"tech-challenge-user-validation/internal/core/usecases"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// Secret should come from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-for-local-dev"
	}

	// Initialize repositories
	userRepo := repositories.NewRDSUserRepository()
	tokenRepo := repositories.NewDynamoTokenRepository()

	// Initialize usecase
	authUseCase := usecases.NewAuthUseCase(userRepo, tokenRepo, jwtSecret)

	// Initialize handler
	authHandler := handlers.NewAuthHandler(authUseCase)

	// Start Lambda
	lambda.Start(authHandler.Handle)
}
