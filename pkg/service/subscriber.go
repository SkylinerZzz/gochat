package service

import (
	log "github.com/sirupsen/logrus"
	"gochat/common"
	"gochat/pkg/queue"
)

type Subscriber struct {
	queue *queue.Queue
}

func NewSubscriber(queue *queue.Queue) *Subscriber {
	return &Subscriber{queue: queue}
}

func (s *Subscriber) Exec(val ...interface{}) error {
	// parameter verification
	if len(val) != 1 {
		log.WithFields(log.Fields{
			"val": val,
		}).Errorf("[Subscriber] wrong number of parameters")
		return ErrInvalidParams
	}
	roomId, ok := val[0].(string)
	if !ok {
		log.WithFields(log.Fields{
			"val": val,
		}).Errorf("[Subscriber] wrong type of parameters")
		return ErrInvalidParams
	}
	if roomId == "" {
		log.WithFields(log.Fields{
			"val": val,
		}).Errorf("[Subscriber] room id can not be empty")
		return ErrInvalidParams
	}
	log.WithFields(log.Fields{
		"roomId": roomId,
	}).Info("[Subscriber] subscribing channel of the room")

	// listen channel
	channel := getChannel(roomId)
	subChan := s.queue.Subscribe(channel)
	go func() {
		for {
			select {
			case msg, ok := <-subChan:
				if !ok {
					log.Info("[Subscriber] channel closed")
					return
				}
				if err := s.process(msg); err != nil {
					log.WithFields(log.Fields{
						"message": msg,
					}).Errorf("[Subscriber] failed to process message")
				}
			}
		}
	}()
	return nil
}

func (s *Subscriber) Name() string {
	return "subscriber"
}

func (s *Subscriber) process(message queue.Message) error {
	return nil
}

func getChannel(roomId string) string {
	return common.PREFIX_CHANNEL + roomId
}
