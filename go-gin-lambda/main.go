package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/jackc/pgx/v5/stdlib" // Driver PostgreSQL
)

var (
	db   *sql.DB
	once sync.Once
)

// initDB garante que a conexão seja aberta apenas uma vez por ciclo de vida do container Lambda
func initDB() *sql.DB {
	once.Do(func() {
		// As credenciais devem vir das variáveis de ambiente definidas no deploy.sh ou docker-compose
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

		var err error
		db, err = sql.Open("pgx", connStr)
		if err != nil {
			log.Fatalf("Erro ao abrir banco: %v", err)
		}

		// Configurações de pool essenciais para Lambda
		db.SetMaxOpenConns(1) // Lambda geralmente processa 1 evento por vez
		db.SetMaxIdleConns(1)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err := db.Ping(); err != nil {
			log.Fatalf("Erro ao conectar no banco: %v", err)
		}
		log.Println("Conexão com o banco estabelecida com sucesso!")
	})
	return db
}

func main() {
	// Inicializa o banco antes de subir o servidor Gin
	initDB()

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		// Exemplo de uso: verificando a saúde do banco na rota
		err := db.Ping()
		if err != nil {
			c.JSON(500, gin.H{"status": "database error", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "pong", "db_status": "connected"})
	})

	// No Lambda com Gin, geralmente usamos um adaptador (como o aws-lambda-go-api-proxy)
	// Mas para rodar localmente ou via binário 'bootstrap' conforme seu deploy.sh:
	r.Run(":8080")
}

//package main
//
//import (
//	"context"
//	"log"
//
//	"github.com/aws/aws-lambda-go/events"
//	"github.com/aws/aws-lambda-go/lambda"
//	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
//	"github.com/gin-gonic/gin"
//)
//
//var ginLambda *ginadapter.GinLambda
//
//func init() {
//	// stdout and stderr are sent to AWS CloudWatch Logs
//	log.Printf("Gin cold start")
//	r := gin.Default()
//	r.GET("/hello", func(c *gin.Context) {
//		c.JSON(200, gin.H{
//			"message": "hello world from lambda",
//		})
//	})
//
//	ginLambda = ginadapter.New(r)
//}
//
//func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	// If no name is provided in the HTTP request body, throw an error
//	return ginLambda.ProxyWithContext(ctx, req)
//}
//
//func main() {
//	lambda.Start(Handler)
//}
