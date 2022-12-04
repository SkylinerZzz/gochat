package task

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gochat/common"
	"gochat/pkg/adapter"
	"gochat/pkg/queue"
	"time"
)

// Dispatcher is in charge of dispatching normal message and handling online and offline message
type Dispatcher struct {
	queue *queue.Queue
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
		err = d.processOnline(wsMessage)
	case common.WsMessageTypeContent:
		err = d.processContent(wsMessage)
	case common.WsMessageTypeImage:
		err = d.processImage(wsMessage)
	case common.WsMessageTypeOffline:
		err = d.processOffline(wsMessage)
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
func (d *Dispatcher) processOnline(data common.WsMessage) error {
	return nil
}

// process content type of WsMessage
func (d *Dispatcher) processContent(data common.WsMessage) error {
	return nil
}

// process image type of WsMessage
func (d *Dispatcher) processImage(data common.WsMessage) error {
	return nil
}

// process offline type of WsMessage
func (d *Dispatcher) processOffline(data common.WsMessage) error {
	return nil
}
