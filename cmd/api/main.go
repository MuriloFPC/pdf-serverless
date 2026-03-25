package main

import (
	"context"
	"log"
	"os"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/domain/interfaces"
	"pdf_serverless/internal/core/service/pdf_service"
	"pdf_serverless/internal/core/service/pdf_service/strategy"
	"pdf_serverless/internal/infra/database"
	"pdf_serverless/internal/infra/queue"
	"pdf_serverless/internal/infra/storage"
	"pdf_serverless/internal/router"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	ctx := context.Background()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	// AWS Configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Initialize Infrastructure
	var jobRepo entities.PDFJobRepository
	var userRepo entities.UserRepository
	var store interfaces.StorageProvider
	var q interfaces.QueueProvider

	if os.Getenv("USE_MEMORY") == "true" {
		jobRepo = database.NewJobMemoryRepository()
		userRepo = database.NewUserMemoryRepository()
		store = storage.NewMemoryStorage()
		q = queue.NewMemoryQueue()
	} else {
		dynamoClient := dynamodb.NewFromConfig(cfg)
		s3Client := s3.NewFromConfig(cfg)
		sqsClient := sqs.NewFromConfig(cfg)

		jobRepo = database.NewJobDynamoRepository(dynamoClient)
		userRepo = database.NewUserDynamoRepository(dynamoClient)
		store = storage.NewS3Storage(s3Client)
		q = queue.NewSQSQueue(sqsClient)
	}

	// Initialize PDF Service & Strategies
	strategies := []strategy.ProcessingStrategy{
		strategy.NewMergeStrategy(store),
		strategy.NewSplitStrategy(store),
		strategy.NewProtectStrategy(store),
		strategy.NewUnprotectStrategy(store),
	}
	pdfService := pdf_service.NewPDFService(jobRepo, store, q, strategies)

	// Initialize Handlers
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecretkey"
	}
	authHandler := router.NewAuthHandler(userRepo, jwtSecret)
	pdfHandler := router.NewPDFHandler(pdfService, store)

	// Start worker in background (only for memory-based demo)
	if os.Getenv("USE_MEMORY") == "true" {
		go func() {
			if memQ, ok := q.(*queue.MemoryQueue); ok {
				for jobID := range memQ.Messages() {
					log.Printf("Worker: Processing job %s", jobID)
					if err := pdfService.ProcessJob(ctx, jobID); err != nil {
						log.Printf("Worker: Error processing job %s: %v", jobID, err)
					} else {
						log.Printf("Worker: Job %s completed", jobID)
					}
				}
			}
		}()
	}

	// Routes
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	pdf := app.Group("/pdf", authHandler.JWTMiddleware())
	pdf.Post("/process", pdfHandler.Process)
	pdf.Get("/status/:id", pdfHandler.GetStatus)
	pdf.Get("/list", pdfHandler.List)

	if os.Getenv("LAMBDA_TASK_ROOT") != "" || os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		adapter := fiberadapter.New(app)
		lambda.Start(adapter.ProxyWithContext)
	} else {
		log.Fatal(app.Listen(":3000"))
	}
}
