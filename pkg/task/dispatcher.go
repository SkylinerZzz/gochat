package task

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gochat/common"
	"gochat/modelv2"
	"gochat/pkg/adapter"
	"gochat/pkg/queue"
	"gochat/pkg/service"
	"time"
)

// Dispatcher is in charge of dispatching normal message and handling online and offline message
type Dispatcher struct {
	queue *queue.Queue
}

func NewDispatcher(queue *queue.Queue) *Dispatcher {
	return &Dispatcher{
		queue: queue,
	}
}

func (d *Dispatcher) Run(message queue.Message) (info adapter.QueueTaskInfo, status adapter.QueueTaskStatus, err error) {
	// record task info
	startTime := time.Now()
	info.Message = message
	info.TaskName = d.Name()
	defer func() {
		info.Duration = time.Since(startTime)
	}()

	log.WithFields(log.Fields{
		"queueName": message.QueueName,
		"data":      message.Data,
	}).Info("[Dispatcher] start to process message")
	// resolve data of message
	wsMessage := common.WsMessage{}
	err = json.Unmarshal([]byte(message.Data), &wsMessage)
	if err != nil {
		log.Errorf("[Dispatcher] failed to unmarshal message, err =%s", err)
		return info, adapter.QueueTaskStatusFailure, err
	}
	// process certain types of WsMessage
	switch wsMessage.Type {
	case common.WsMessageTypeOnline:
		err = d.processOnline(message.Data)
	case common.WsMessageTypeContent:
		err = d.processContent(message.Data)
	case common.WsMessageTypeImage:
		err = d.processImage(message.Data)
	case common.WsMessageTypeOffline:
		err = d.processOffline(message.Data)
	default:
		log.WithFields(log.Fields{
			"type": wsMessage.Type,
		}).Error("[Dispatcher] unknown type of WsMessage")
		err = ErrUnknownWsMessageType
	}
	if err != nil {
		log.Errorf("[Dispatcher] failed to process WsMessage, err = %s", err)
		return info, adapter.QueueTaskStatusFailure, err
	}

	return info, adapter.QueueTaskStatusSuccess, nil
}

func (d *Dispatcher) Name() string {
	return "dispatcher"
}

// process online type of WsMessage
func (d *Dispatcher) processOnline(data string) error {
	// get room id
	message := modelv2.Message{}
	wsMessage := common.WsMessage{Data: &message}
	err := json.Unmarshal([]byte(data), &wsMessage)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"roomId": message.RoomId,
		"userId": message.UserId,
	}).Info("[Dispatcher] process online message")

	// broadcast directly
	bd := service.NewBroadcaster()
	err = bd.Exec(wsMessage, message.RoomId)
	if err != nil {
		log.Errorf("[Dispatcher] failed to broadcast online message, err = %s", err)
		return err
	}
	return nil
}

// process content type of WsMessage
func (d *Dispatcher) processContent(data string) error {
	message := queue.Message{
		QueueName: common.DATABUS_CONTENT_HANDLER,
		Data:      data,
	}
	d.queue.SendMessage(common.DATABUS_CONTENT_HANDLER, message)
	return nil
}

// process image type of WsMessage
func (d *Dispatcher) processImage(data string) error {
	return nil
}

// process offline type of WsMessage
func (d *Dispatcher) processOffline(data string) error {
	return nil
}
