package entities

import (
	"context"
	"time"
)

type User struct {
	ID           string    `json:"id" dynamodbav:"id"`
	Email        string    `json:"email" dynamodbav:"email"`
	PasswordHash string    `json:"-" dynamodbav:"password_hash"`
	CreatedAt    time.Time `json:"created_at" dynamodbav:"created_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}
