package queue

import (
	"context"
	"fmt"
)

type MemoryQueue struct {
	messages chan string
}

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		messages: make(chan string, 100),
	}
}

func (q *MemoryQueue) Publish(ctx context.Context, jobID string) error {
	select {
	case q.messages <- jobID:
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

func (q *MemoryQueue) Messages() <-chan string {
	return q.messages
}
