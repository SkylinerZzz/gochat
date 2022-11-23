package queue

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// Queue encapsulates Node interface
type Queue struct {
	node Node
}

func NewQueue(config map[string]string) (*Queue, error) {
	node, err := NewRedisNode(config)
	if err != nil {
		log.Errorf("[Queue] failed to init redis node, err = %s", err)
		return nil, err
	}
	return &Queue{node: node}, nil
}

func (q *Queue) ReceiveMessage(queueName string, timeout time.Duration) (Message, error) {
	return q.node.ReceiveMessage(queueName, timeout)
}

func (q *Queue) Subscribe(channel string) <-chan Message {
	return q.node.Subscribe(channel)
}

func (q *Queue) SendMessage(queueName string, message Message) error {
	return q.node.SendMessage(queueName, message)
}

func (q *Queue) Publish(channel string, message Message) error {
	return q.node.Publish(channel, message)
}
