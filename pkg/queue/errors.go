package queue

import "errors"

var (
	ErrQueueEmpty = errors.New("redisqueue: queue is empty")
)
