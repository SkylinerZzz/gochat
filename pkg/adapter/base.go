package adapter

import (
	"context"
	"gochat/pkg/queue"
)

// QueueTask interface
type QueueTask interface {
	Run(ctx context.Context, message queue.Message) error
	Name() string
}

// Handler interface
type Handler interface {
	Handle(err error)
	Name() string
}
