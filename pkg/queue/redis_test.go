package queue_test

import (
	"context"
	"fmt"
	"gochat/pkg/queue"
	"gochat/util"
	"testing"
	"time"
)

const testQueue = "gochat:test:queue"

func TestQueue(t *testing.T) {
	util.Init("../../config")
	q := util.RedisQueue
	req := queue.Message{
		Data:      "hello",
		QueueName: testQueue,
	}
	err := q.SendMessage(testQueue, req)
	if err != nil {
		t.Error(err)
	}
	resp, err := q.ReceiveMessage(testQueue)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}

func TestPubsub(t *testing.T) {
	util.Init("../../config")
	q := util.RedisQueue
	channel := "gochat:test:pubsub"
	subChan := q.Subscribe(channel)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("unsubscribe")
				return
			case msg := <-subChan:
				fmt.Println(msg)
			}
		}
	}()
	time.Sleep(10 * time.Second)
}
