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

type ContentHandler struct {
	queue *queue.Queue
}

func (ch *ContentHandler) Run(message queue.Message) (info adapter.QueueTaskInfo, status adapter.QueueTaskStatus, err error) {
	// record task info
	startTime := time.Now()
	info.Message = message
	info.TaskName = ch.Name()
	defer func() {
		info.Duration = time.Since(startTime)
	}()

	data := modelv2.Message{}
	wsMessage := common.WsMessage{Data: &data}
	err = json.Unmarshal([]byte(message.Data), &wsMessage)
	if err != nil {
		log.WithFields(log.Fields{
			"queueName": message.QueueName,
			"data":      message.Data,
		}).Errorf("[ContentHandler] failed to unmarshal message, err = %s", err)
		return info, adapter.QueueTaskStatusFailure, err
	}

	// save message
	err = modelv2.SaveMessage(data)
	if err != nil {
		log.WithFields(log.Fields{
			"userId":  data.UserId,
			"roomId":  data.RoomId,
			"content": data.Content,
		}).Errorf("[ContentHandler] failed to save message, err = %s", err)
		return info, adapter.QueueTaskStatusFailure, err
	}

	// broadcast message
	bd := service.NewBroadcaster()
	err = bd.Exec(wsMessage, data.RoomId)
	if err != nil {
		log.WithFields(log.Fields{
			"wsMessage": wsMessage,
			"roomId":    data.RoomId,
		}).Errorf("[ContentHandler] failed to broadcast message, err = %s", err)
		return info, adapter.QueueTaskStatusFailure, err
	}

	return info, adapter.QueueTaskStatusSuccess, nil
}
func (ch *ContentHandler) Name() string {
	return "content handler"
}
