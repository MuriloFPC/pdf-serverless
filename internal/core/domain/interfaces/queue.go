package interfaces

import (
	"context"
)

type QueueProvider interface {
	Publish(ctx context.Context, jobID string) error
}
