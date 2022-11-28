package adapter

import (
	"context"
	"gochat/pkg/queue"
)

type QueueTask interface {
	Run(ctx context.Context, message queue.Message) error
	Name() string
}
