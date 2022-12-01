package adapter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	success int // task completion count
	failure int // task incomplete count
	slow    int // task timeout count
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Handle(info QueueTaskInfo, status QueueTaskStatus, err error) {
	switch status {
	case QueueTaskStatusSuccess: // err is nil
		l.success++
		log.WithFields(log.Fields{
			"queueName": info.Message.QueueName,
			"taskName":  info.TaskName,
			"data":      info.Message.Data,
		}).Info("[QueueTaskAdapter] process message successfully")
	case QueueTaskStatusFailure: // err not nil
		l.failure++
		log.WithFields(log.Fields{
			"queueName": info.Message.QueueName,
			"taskName":  info.TaskName,
			"data":      info.Message.Data,
		}).Errorf("[QueueTaskAdapter] failed to process message, err = %s", err)
	case QueueTaskStatusTimeout: // err is nil but timeout
		l.slow++
		log.WithFields(log.Fields{
			"queueName": info.Message.QueueName,
			"taskName":  info.TaskName,
			"data":      info.Message.Data,
		}).Info("[QueueTaskAdapter] process message successfully but slow")
	}
}

func (l *Logger) Name() string {
	return "Logger"
}

func (l *Logger) Log() {
	rate := float32(l.success) / float32(l.success+l.failure+l.slow) * 100
	fmt.Printf("the number of successful task: %d\n"+
		"the number of failed task: %d\n"+
		"the number of slow task: %d\n"+
		"success rate: %.1f%%\n", l.success, l.failure, l.slow, rate)
}
