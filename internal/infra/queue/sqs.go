package queue

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSQueue struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSQueue(client *sqs.Client) *SQSQueue {
	queueURL := os.Getenv("SQS_QUEUE_URL")
	// Se não houver URL nas variáveis de ambiente, tentaremos usar o nome padrão
	// Mas o ideal para SQS é ter a URL completa.
	return &SQSQueue{
		client:   client,
		queueURL: queueURL,
	}
}

func (q *SQSQueue) Publish(ctx context.Context, jobID string) error {
	if q.queueURL == "" {
		// Tenta obter a URL pelo nome se não foi configurada
		out, err := q.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: aws.String("pdf-serverless-queue"),
		})
		if err != nil {
			return fmt.Errorf("failed to get SQS queue URL: %w", err)
		}
		q.queueURL = *out.QueueUrl
	}

	_, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(q.queueURL),
		MessageBody: aws.String(jobID),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	return nil
}
