package adapter

import (
	"context"
	"fmt"
	"gochat/pkg/queue"
	"gochat/util"
	"strconv"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, _ := context.WithTimeout(ctx, 5*time.Second)
	go doContext(c)
	time.Sleep(3 * time.Second)
}

func doContext(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			return
		default:
			fmt.Println("run")
			time.Sleep(1 * time.Second)
		}
	}
}

type service struct{}

func (service) Run(ctx context.Context, message queue.Message) error {
	fmt.Println(message)
	return nil
}

func (service) Name() string {
	return "service"
}

func TestAdapter(t *testing.T) {
	util.Init("../../config")
	queueName := "gochat:test:adapter"
	task := service{}
	go func() {
		for i := 0; i < 100000; i++ {
			msg := queue.Message{
				Data:      strconv.Itoa(i),
				QueueName: queueName,
			}
			util.RedisQueue.SendMessage(queueName, msg)
		}
	}()
	logger := NewLogger("logger")
	adapter := NewQueueTaskAdapter(task, util.RedisQueue, queueName, 5*time.Second, 100, logger)
	go adapter.Start()
	time.Sleep(1 * time.Minute)
	adapter.Terminate()
	fmt.Println("success:", logger.success)
}
