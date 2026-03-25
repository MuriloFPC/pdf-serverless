package main

import (
	"context"
	"log"
	"pdf_serverless/internal/core/service/pdf_service"
	"pdf_serverless/internal/core/service/pdf_service/strategy"
	"pdf_serverless/internal/infra/database"
	"pdf_serverless/internal/infra/storage"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	ctx := context.Background()

	// AWS Configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)

	// Initialize Infrastructure
	jobRepo := database.NewJobDynamoRepository(dynamoClient)
	store := storage.NewS3Storage(s3Client)

	// Initialize PDF Service & Strategies
	strategies := []strategy.ProcessingStrategy{
		strategy.NewMergeStrategy(store),
		strategy.NewSplitStrategy(store),
		strategy.NewProtectStrategy(store),
		strategy.NewUnprotectStrategy(store),
	}
	pdfService := pdf_service.NewPDFService(jobRepo, store, nil, strategies)

	handler := func(ctx context.Context, sqsEvent events.SQSEvent) error {
		for _, message := range sqsEvent.Records {
			jobID := message.Body
			log.Printf("Processing job: %s", jobID)

			if err := pdfService.ProcessJob(ctx, jobID); err != nil {
				log.Printf("Error processing job %s: %v", jobID, err)
				// Returning error will make the message visible again in the queue (if visibility timeout allows)
				return err
			}
			log.Printf("Job %s completed", jobID)
		}
		return nil
	}

	lambda.Start(handler)
}
