package interfaces

import "context"

type EmailService interface {
	Send(ctx context.Context, to, subject, body string) error
}
