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

type UserDynamoRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewUserDynamoRepository(client *dynamodb.Client) *UserDynamoRepository {
	tableName := os.Getenv("USER_TABLE")
	if tableName == "" {
		tableName = "pdf-serverless-user"
	}
	return &UserDynamoRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *UserDynamoRepository) Create(ctx context.Context, user *entities.User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
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

func (r *UserDynamoRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("EmailIndex"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}

	if len(out.Items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user entities.User
	err = attributevalue.UnmarshalMap(out.Items[0], &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

func (r *UserDynamoRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if out.Item == nil {
		return nil, fmt.Errorf("user not found")
	}

	var user entities.User
	err = attributevalue.UnmarshalMap(out.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}
