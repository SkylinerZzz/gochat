package service

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gochat/common"
	"gochat/pkg/queue"
	"gochat/util"
	"sync"
)

type Broadcaster struct{}

var (
	broadcasterService *Broadcaster
	broadcasterOnce    sync.Once
)

func NewBroadcaster() *Broadcaster {
	broadcasterOnce.Do(func() {
		broadcasterService = &Broadcaster{}
	})
	return broadcasterService
}

func (b *Broadcaster) Exec(val ...interface{}) error {
	// parameter verification
	// val[0] WsMessage
	// val[1] RoomId
	if len(val) != 2 {
		log.WithFields(log.Fields{
			"val": val,
		}).Errorf("[Broadcaster] wrong number of parameters")
		return ErrInvalidParams
	}
	wsMessage, ok := val[0].(common.WsMessage)
	if !ok {
		log.WithFields(log.Fields{
			"val[0]": val[0],
		}).Errorf("[Broadcaster] wrong type of parameters")
		return ErrInvalidParams
	}
	roomId, ok := val[1].(string)
	if !ok {
		log.WithFields(log.Fields{
			"val[1]": val[1],
		}).Errorf("[Broadcaster] wrong type of parameters")
		return ErrInvalidParams
	}
	log.WithFields(log.Fields{
		"wsMessage": wsMessage,
		"roomId":    roomId,
	}).Info("[Broadcaster] broadcasting")

	rd := util.RedisPool.Get()
	defer rd.Close()
	// traverse user list of given room
	userList, err := redis.IntMap(rd.Do("hgetall", getUserListKey(roomId)))
	if err != nil {
		log.WithFields(log.Fields{
			"userListKey": getUserListKey(roomId),
		}).Errorf("[Broadcaster] failed to load user list, err = %s", err)
		return err
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		log.WithFields(log.Fields{
			"wsMessage": wsMessage,
			"roomId":    roomId,
		}).Errorf("[Broadcaster] failed to marshal ws message, err = %s", err)
		return err
	}
	for userId, status := range userList {
		if status == common.UserStatusOffline {
			continue
		}
		// search local client map
		v, ok := common.ClientMap[roomId].Load(userId)
		common.ClientMap[roomId].Range(func(key any, value any) bool {
			log.WithFields(log.Fields{
				"userId": userId,
			}).Infof("[Debug] ??%v", key)
			return false
		})
		if !ok {
			// publish ws message to let other server handle this message that has corresponding client map
			// fill PubSubMessage
			pub := common.PubSubMessage{
				UserId: userId,
				RoomId: roomId,
				Data:   string(data),
			}
			pubData, err := json.Marshal(pub)
			if err != nil {
				log.WithFields(log.Fields{
					"targetUser": userId,
					"data":       pub.Data,
				}).Errorf("[Broadcaster] failed to marshal pub message, err = %s", err)
				return err
			}

			message := queue.Message{
				Data: string(pubData),
			}
			util.RedisQueue.Publish(getChannel(roomId), message)
			continue
		}
		wsClient := v.(common.WsClient)
		err = wsClient.Conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Errorf("[Broadcaster] failed to write ws message, err = %s", err)
		}
	}
	return nil
}

func (b *Broadcaster) Name() string {
	return "broadcaster"
}
