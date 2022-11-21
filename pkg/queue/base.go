package queue

import "time"

// Message definition
type Message struct {
	Data      string `json:"data"`       // raw data
	QueueName string `json:"queue_name"` // owner
}

// Queue interface
type Queue interface {
	Consumer
	Provider
}

// Consumer interface
type Consumer interface {
	ReceiveMessage(queueName string, timeout time.Duration) (Message, error)
	Subscribe(channel string) <-chan Message
}

// Provider interface
type Provider interface {
	SendMessage(queueName string, message Message) error
	Publish(channel string, message Message) error
}
