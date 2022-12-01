package queue

import "errors"

var (
	ErrQueueEmpty = errors.New("RedisQueue: queue is empty")
)
