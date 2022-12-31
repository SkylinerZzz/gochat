package service

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gochat/common"
	"gochat/pkg/queue"
	"gochat/util"
	"sync"
	"time"
)

// Entry is in charge of initializing map and sending to dispatcher
type Entry struct{}

var (
	entryService *Entry
	entryOnce    sync.Once
)

func NewEntry() *Entry {
	entryOnce.Do(func() {
		entryService = &Entry{}
	})
	return entryService
}

func (e *Entry) Exec(val ...interface{}) error {
	// parameter verification
	// val[0]: websocket.Conn
	// val[1]: RoomId
	// val[2]: UserId
	if len(val) != 3 {
		log.WithFields(log.Fields{
			"val": val,
		}).Errorf("[Entry] wrong number of parameters")
		return ErrInvalidParams
	}
	conn, ok := val[0].(*websocket.Conn)
	if !ok {
		log.WithFields(log.Fields{
			"val[0]": val[0],
		}).Errorf("[Entry] wrong type of parameters")
		return ErrInvalidParams
	}
	roomId, ok := val[1].(string)
	if !ok {
		log.WithFields(log.Fields{
			"val[1]": val[1],
		}).Errorf("[Entry] wrong type of parameters")
		return ErrInvalidParams
	}
	userId, ok := val[2].(string)
	if !ok {
		log.WithFields(log.Fields{
			"val[2]": val[2],
		}).Errorf("[Entry] wrong type of parameters")
		return ErrInvalidParams
	}
	log.WithFields(log.Fields{
		"roomId": roomId,
		"userId": userId,
	}).Info("[Entry] prepare to receive ws message")

	wsClient := common.WsClient{
		Conn:   conn,
		RoomId: roomId,
		UserId: userId,
	}
	if common.ClientMap[roomId] == nil {
		common.ClientMapMutex.Lock()
		if common.ClientMap[roomId] == nil {
			common.ClientMap[roomId] = &sync.Map{}
			sub := NewSubscriber()
			err := sub.Exec(roomId)
			if err != nil {
				log.WithFields(log.Fields{
					"roomId": roomId,
				}).Errorf("[Entry] failed to subscribe channel, err = %s", err)
				return err
			}
		}
		common.ClientMapMutex.Unlock()
	}

	// online
	checkOldClientClosed(roomId, userId)

	common.ClientMap[roomId].Store(userId, wsClient)
	rd := util.RedisPool.Get()
	defer rd.Close()
	_, err := rd.Do("hset", getUserListKey(roomId), userId, common.UserStatusOnline)
	if err != nil {
		log.WithFields(log.Fields{
			"roomId": roomId,
			"userId": userId,
		}).Errorf("[Entry] failed to set user online status, err = %s", err)
		return err
	}

	go read(wsClient)

	return nil
}

func (e *Entry) Name() string {
	return "entry"
}

func read(ws common.WsClient) {
	defer ws.Conn.Close()
	for {
		t, data, err := ws.Conn.ReadMessage()
		if err != nil {
			if err.Error() == ErrWsConnClosed.Error() {
				log.Info("[Entry] websocket connection closed")
				return
			}
			log.Errorf("[Entry] failed to receive ws message, err = %s", err)
			return
		}
		switch t {
		// offline
		case websocket.CloseMessage:
			// pass
			return
		default:
			message := queue.Message{
				QueueName: common.DATABUS_DISPATCHER,
				Data:      string(data),
			}
			err := util.RedisQueue.SendMessage(common.DATABUS_DISPATCHER, message)
			if err != nil {
				log.WithFields(log.Fields{
					"data":      message.Data,
					"queueName": message.QueueName,
				}).Errorf("[Entry] failed to send message to dispatcher, err = %s", err)
				return
			}
			log.WithFields(log.Fields{
				"data":      message.Data,
				"queueName": message.QueueName,
			}).Info("[Entry] send message to dispatcher successfully")
		}
	}
}

func getUserListKey(roomId string) string {
	return common.PREFIX_USER_LIST + roomId
}

func checkOldClientClosed(roomId, userId string) {
	// confirm old client has been removed
	maxTry := 10
	for curTry := 0; curTry < maxTry; curTry++ {
		if _, ok := common.ClientMap[roomId].Load(userId); ok {
			time.Sleep(100 * time.Millisecond)
		} else {
			return
		}
	}
	// force remove old client
	common.ClientMap[roomId].Delete(userId)
	log.WithFields(log.Fields{
		"roomId": roomId,
		"userId": userId,
	}).Info("[Entry] retry times exceeded, remove old client forcibly")
}
