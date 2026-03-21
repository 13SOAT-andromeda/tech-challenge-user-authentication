package main

import (
	"context"
	"fmt"
	"log"
	"os"
	services "tech-challenge-user-validation/internal/service"
	"time"

	"tech-challenge-user-validation/internal/adapters/database/model/address"
	"tech-challenge-user-validation/internal/adapters/database/model/user"
	"tech-challenge-user-validation/internal/adapters/database/repositories"
	"tech-challenge-user-validation/internal/adapters/handlers"
	jwtadapter "tech-challenge-user-validation/internal/adapters/services/jwt"
	"tech-challenge-user-validation/internal/core/usecases"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		panic("JWT_SECRET is not set")
	}

	// 2. Postgres Connection (GORM)
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	err = db.AutoMigrate(&user.Model{}, &address.Model{})
	if err != nil {
		log.Printf("failed to auto migrate models: %v", err)
	}

	// 3. DynamoDB Connection
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "user-auth-tokens"
	}

	// 4. Dependency Injection
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(dynamoClient, tableName)

	jwtSvc := jwtadapter.NewService(jwtSecret, 15*time.Minute, 24*time.Hour)
	sessionSvc := services.NewSessionService(sessionRepo)
	authUseCase := usecases.NewAuthUseCase(userRepo, nil, sessionSvc, jwtSvc, jwtSecret)

	authHandler := handlers.NewAuthHandler(authUseCase)

	// 5. Start Lambda
	log.Println("Starting Lambda...")
	lambda.Start(authHandler.Handle)
}
