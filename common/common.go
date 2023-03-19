package common

import (
	"github.com/gorilla/websocket"
	"sync"
)

// WsClient definition
type WsClient struct {
	Conn   *websocket.Conn // websocket connection
	UserId string          // owner id
	RoomId string          // room id
}

var (
	ClientMapMutex sync.RWMutex
	ClientMap      = map[string]*sync.Map{} // store WsClient
)

// type of WsMessage
const (
	WsMessageTypeOnline = iota
	WsMessageTypeContent
	WsMessageTypeImage
	WsMessageTypeOffline
)

// WsMessage describe websocket raw data
type WsMessage struct {
	Type int         `json:"type"` // WsMessage type
	Data interface{} `json:"data"` // WsMessage data
}

const (
	DATABUS_DISPATCHER      = "gochat:test:dispatcher"      // consumer queue of dispatcher
	DATABUS_CONTENT_HANDLER = "gochat:test:content_handler" // consumer queue of content handler
	DATABUS_IMAGE_HANDLER   = "gochat:test:image_handler"   // consumer queue of image handler
	PREFIX_CHANNEL          = "gochat:test:channel:room_"   // prefix of channel
	PREFIX_USER_LIST        = "gochat:test:list:room_"      // prefix of user list in each room, recording whether a user is offline or online
	PREFIX_PRIVATE_ROOM     = "gochat:test:private:room_"   // prefix of private chat, recording corresponding room id
)

// status of user
const (
	UserStatusOnline = iota
	UserStatusOffline
)

// PubSubMessage wraps WsMessage, sending to given user
type PubSubMessage struct {
	UserId string `json:"user_id"` // target user
	RoomId string `json:"room_id"` // room id
	Data   string `json:"data"`    // WsMessage
}
