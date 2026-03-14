package main

import (
	"context"
	"log"
	"pdf_serverless/internal/core/service/pdf_service"
	"pdf_serverless/internal/core/service/pdf_service/strategy"
	"pdf_serverless/internal/infra/database"
	"pdf_serverless/internal/infra/queue"
	"pdf_serverless/internal/infra/storage"
	"pdf_serverless/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	// Initialize Infrastructure
	jobRepo := database.NewJobMemoryRepository()
	userRepo := database.NewUserMemoryRepository()
	store := storage.NewMemoryStorage()
	q := queue.NewMemoryQueue()

	// Initialize PDF Service & Strategies
	strategies := []strategy.ProcessingStrategy{
		strategy.NewMergeStrategy(store),
		strategy.NewSplitStrategy(store),
		strategy.NewProtectStrategy(store),
		strategy.NewUnprotectStrategy(store),
	}
	pdfService := pdf_service.NewPDFService(jobRepo, store, q, strategies)

	// Initialize Handlers
	jwtSecret := "supersecretkey" // In production, use env variable
	authHandler := router.NewAuthHandler(userRepo, jwtSecret)
	pdfHandler := router.NewPDFHandler(pdfService, store)

	// Start worker in background (for memory-based demo)
	go func() {
		ctx := context.Background()
		for jobID := range q.Messages() {
			log.Printf("Worker: Processing job %s", jobID)
			if err := pdfService.ProcessJob(ctx, jobID); err != nil {
				log.Printf("Worker: Error processing job %s: %v", jobID, err)
			} else {
				log.Printf("Worker: Job %s completed", jobID)
			}
		}
	}()

	// Routes
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	pdf := app.Group("/pdf", authHandler.JWTMiddleware())
	pdf.Post("/process", pdfHandler.Process)
	pdf.Get("/status/:id", pdfHandler.GetStatus)
	pdf.Get("/list", pdfHandler.List)

	log.Fatal(app.Listen(":3000"))
}
