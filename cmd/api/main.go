package main

import (
	"context"
	"log"
	"os"
	_ "pdf_serverless/docs"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/domain/interfaces"
	"pdf_serverless/internal/core/service/pdf_service"
	"pdf_serverless/internal/core/service/pdf_service/strategy"
	"pdf_serverless/internal/infra/database"
	"pdf_serverless/internal/infra/email"
	"pdf_serverless/internal/infra/queue"
	"pdf_serverless/internal/infra/storage"
	"pdf_serverless/internal/router"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title PDF Serverless API
// @version 1.0
// @description This is a serverless PDF processing API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Handle OPTIONS request specifically for Lambda compatibility if needed
	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

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
		strategy.NewOptimizeStrategy(store),
	}
	pdfService := pdf_service.NewPDFService(jobRepo, store, q, strategies)

	// Initialize Handlers
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecretkey"
	}
	emailService := email.NewNoOpEmailService()
	authHandler := router.NewAuthHandler(userRepo, emailService, jwtSecret)
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
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "PDF Serverless API is running",
		})
	})

	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	pdf := app.Group("/pdf", authHandler.JWTMiddleware())
	pdf.Post("/process", pdfHandler.Process)
	pdf.Get("/presigned-url/:id", pdfHandler.GetPresignedURL)
	pdf.Post("/complete-upload/:id", pdfHandler.CompleteUpload)
	pdf.Get("/status/:id", pdfHandler.GetStatus)
	pdf.Get("/download/:id", pdfHandler.GetDownloadURL)
	pdf.Get("/list", pdfHandler.List)
	pdf.Delete("/:id", pdfHandler.Delete)

	if os.Getenv("LAMBDA_TASK_ROOT") != "" || os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		adapter := fiberadapter.New(app)
		lambda.Start(adapter.ProxyWithContextV2)
	} else {
		log.Fatal(app.Listen(":3000"))
	}
}
