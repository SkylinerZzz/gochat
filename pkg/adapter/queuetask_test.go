package adapter_test

import (
	"gochat/pkg/adapter"
	"gochat/pkg/queue"
	"gochat/util"
	"strconv"
	"sync"
	"testing"
	"time"
)

type TestTask struct {
	queue     *queue.Queue
	queueName string
}

var testTask *TestTask
var testTaskOnce sync.Once

func NewTestTask(queue *queue.Queue, queueName string) *TestTask {
	testTaskOnce.Do(func() {
		testTask = &TestTask{
			queue:     queue,
			queueName: queueName,
		}
	})
	return testTask
}

func (t *TestTask) Run(message queue.Message) (info adapter.QueueTaskInfo, status adapter.QueueTaskStatus, err error) {
	start := time.Now()
	defer func() {
		info.Duration = time.Since(start)
	}()

	info.TaskName = t.Name()
	info.Message = message

	num, _ := strconv.Atoi(message.Data)
	if num%2 == 0 {
		time.Sleep(2 * time.Second)
	}
	if err := t.queue.SendMessage(t.queueName, message); err != nil {
		return info, adapter.QueueTaskStatusFailure, err
	}
	return info, adapter.QueueTaskStatusSuccess, nil
}

func (t *TestTask) Name() string {
	return "TestTask"
}

func TestAdapter(t *testing.T) {
	util.Init("../../config")
	queueName := "gochat:test:adapter-input"
	task := NewTestTask(util.RedisQueue, "gochat:test:adapter-output")
	go func() {
		for i := 0; i < 1000; i++ {
			msg := queue.Message{
				Data:      strconv.Itoa(i),
				QueueName: queueName,
			}
			util.RedisQueue.SendMessage(queueName, msg)
		}
	}()

	logger := adapter.NewLogger()
	adapter := adapter.NewQueueTaskAdapter(task, util.RedisQueue, queueName, 1*time.Second, 100, logger)
	go adapter.Start()
	time.Sleep(20 * time.Second)
	adapter.Terminate()
	logger.Log()
}
