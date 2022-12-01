package adapter

import (
	"gochat/pkg/queue"
	"time"
)

// QueueTaskStatus describes processing status of a given task
type QueueTaskStatus int

const (
	QueueTaskStatusFailure QueueTaskStatus = iota
	QueueTaskStatusSuccess
	QueueTaskStatusTimeout
)

// QueueTaskInfo describes detail info of a given task about its message and task name
type QueueTaskInfo struct {
	Message  queue.Message
	TaskName string
	Duration time.Duration
}

// QueueTask interface
type QueueTask interface {
	Run(message queue.Message) (QueueTaskInfo, QueueTaskStatus, error)
	Name() string
}

// Handler interface
type Handler interface {
	Handle(info QueueTaskInfo, status QueueTaskStatus, err error)
	Name() string
}
