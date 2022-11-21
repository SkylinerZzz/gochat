package wrapper

import "gochat/pkg/queue"

type QueueTask interface {
	Run(message queue.Message) error
}
