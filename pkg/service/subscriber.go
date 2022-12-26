package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gochat/common"
	"gochat/modelv2"
	"gochat/pkg/queue"
	"gochat/util"
)

// Subscriber works in distributed environment
type Subscriber struct{}

func NewSubscriber() *Subscriber {
	return &Subscriber{}
}

func (s *Subscriber) Exec(val ...interface{}) error {
	// parameter verification
	// val[0]: room id
	if len(val) != 1 {
		log.WithFields(log.Fields{
			"val": val,
		}).Errorf("[Subscriber] wrong number of parameters")
		return ErrInvalidParams
	}
	roomId, ok := val[0].(string)
	if !ok {
		log.WithFields(log.Fields{
			"val[0]": val[0],
		}).Errorf("[Subscriber] wrong type of parameters")
		return ErrInvalidParams
	}
	log.WithFields(log.Fields{
		"roomId": roomId,
	}).Info("[Subscriber] subscribing channel of the room")

	// listen channel
	channel := getChannel(roomId)
	subChan := util.RedisQueue.Subscribe(channel)
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
	data := modelv2.Message{}
	wsMessage := common.WsMessage{Data: &data}
	err := json.Unmarshal([]byte(message.Data), &wsMessage)
	if err != nil {
		log.WithFields(log.Fields{
			"queueName": message.QueueName,
			"data":      message.Data,
		}).Errorf("[Subscriber] failed to unmarshal message, err = %s", err)
		return err
	}

	//  search local client map
	v, ok := common.ClientMap[data.RoomId].Load(data.UserId)
	if !ok {
		return nil
	}
	wsClient := v.(common.WsClient)
	log.WithFields(log.Fields{
		"userId": wsClient.UserId,
		"roomId": wsClient.RoomId,
	}).Info("[Subscriber] selected to process message")
	err = wsClient.Conn.WriteMessage(websocket.TextMessage, []byte(message.Data))
	// record error
	if err != nil {
		log.Errorf("[Subscriber] failed to write ws message, err = %s", err)
	}
	return nil
}

func getChannel(roomId string) string {
	return common.PREFIX_CHANNEL + roomId
}
