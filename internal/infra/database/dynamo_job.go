package database

import (
	"context"
	"fmt"
	"os"
	"pdf_serverless/internal/core/domain/entities"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type JobDynamoRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewJobDynamoRepository(client *dynamodb.Client) *JobDynamoRepository {
	tableName := os.Getenv("JOB_TABLE")
	if tableName == "" {
		tableName = "pdf-serverless-job"
	}
	return &JobDynamoRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *JobDynamoRepository) Create(ctx context.Context, job *entities.PDFJob) error {
	item, err := attributevalue.MarshalMap(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

func (r *JobDynamoRepository) GetByID(ctx context.Context, id string) (*entities.PDFJob, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"job_id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	if out.Item == nil {
		return nil, fmt.Errorf("job not found")
	}

	var job entities.PDFJob
	err = attributevalue.UnmarshalMap(out.Item, &job)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

func (r *JobDynamoRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.PDFJob, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("UserIDIndex"),
		KeyConditionExpression: aws.String("user_id = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs by user_id: %w", err)
	}

	var jobs []*entities.PDFJob
	err = attributevalue.UnmarshalListOfMaps(out.Items, &jobs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal jobs: %w", err)
	}

	return jobs, nil
}

func (r *JobDynamoRepository) Update(ctx context.Context, job *entities.PDFJob) error {
	// Reutiliza Create (PutItem sobrescreve se existir a mesma chave primária)
	return r.Create(ctx, job)
}
