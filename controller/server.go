package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gochat/model"
	"strconv"
	"sync"
)

type Client struct {
	Conn     *websocket.Conn
	Username string
	RoomId   string
}
type Message struct {
	MsgType int         `json:"msgType"`
	Data    interface{} `json:"data"`
}

var (
	mutex   = sync.Mutex{}                     // ensure rooms and users concurrency security
	once    = sync.RWMutex{}                   // ensure users[string] should be initialized once
	rooms   = make(map[string][]Client)        // map clients to the room
	users   = make(map[string]map[string]bool) //user mapping
	enter   = make(chan Client, 10)
	leave   = make(chan Client, 10)
	message = make(chan Message, 100)
)

// message type
const (
	_ = iota
	msgTypeOnline
	msgTypeSend
	msgTypeOffline
)

func Run(c *gin.Context) {
	ws, _ := (&websocket.Upgrader{}).Upgrade(c.Writer, c.Request, nil)
	done := make(chan struct{})
	go read(ws, done)
	go write(done)
}

func read(conn *websocket.Conn, done chan struct{}) {
	defer conn.Close()
	defer close(done)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Info("websocket message read error:", err)
			return
		}
		var clientMsg Message
		json.Unmarshal(msg, &clientMsg)
		switch clientMsg.MsgType {
		case msgTypeOnline:
			// user join the room
			enter <- Client{
				Conn:     conn,
				Username: clientMsg.Data.(map[string]interface{})["username"].(string),
				RoomId:   clientMsg.Data.(map[string]interface{})["roomId"].(string),
			}
			message <- clientMsg
		case msgTypeSend:
			// user send message
			strUserId := clientMsg.Data.(map[string]interface{})["userId"].(string)
			strRoomId := clientMsg.Data.(map[string]interface{})["roomId"].(string)
			userId, _ := strconv.ParseUint(strUserId, 10, 32)
			roomId, _ := strconv.ParseUint(strRoomId, 10, 32)
			data := map[string]interface{}{
				"userId":  uint(userId),
				"roomId":  uint(roomId),
				"content": clientMsg.Data.(map[string]interface{})["content"].(string),
			}
			model.SaveContent(data)
			message <- clientMsg
		case msgTypeOffline:
			leave <- Client{
				Conn:     conn,
				Username: clientMsg.Data.(map[string]interface{})["username"].(string),
				RoomId:   clientMsg.Data.(map[string]interface{})["roomId"].(string),
			}
			// do not send leaving message to the channel
		}
	}
}

func write(done chan struct{}) {
	for {
		select {
		case e := <-enter:
			// initial users mapping once
			once.RLock()
			if users[e.RoomId] != nil {
				entering(e)
			}
			once.RUnlock()
			once.Lock()
			if users[e.RoomId] == nil {
				users[e.RoomId] = make(map[string]bool)
			}
			entering(e)
			once.Unlock()
		case l := <-leave:
			leaving(l)
		case msg := <-message:
			logrus.Info("broadcasting ...")
			roomId := msg.Data.(map[string]interface{})["roomId"].(string)
			clients := rooms[roomId]
			for _, c := range clients {
				err := c.Conn.WriteMessage(websocket.TextMessage, formatMessage(msg))
				if err != nil {
					logrus.Warn(err)
				}
			}
		case <-done:
			return
		}
	}
}
func entering(c Client) {
	// update users mapping while entering
	mutex.Lock()
	if _, ok := users[c.RoomId][c.Username]; !ok {
		rooms[c.RoomId] = append(rooms[c.RoomId], c)
		users[c.RoomId][c.Username] = true
	}
	logrus.WithFields(logrus.Fields{
		"roomId":   c.RoomId,
		"roomSize": len(rooms[c.RoomId]),
		"username": c.Username,
	}).Info("an user enter into the room")
	mutex.Unlock()
}
func leaving(c Client) {
	// update users mapping while leaving
	mutex.Lock()
	if _, ok := users[c.RoomId][c.Username]; ok {
		// delete client mapping
		delete(users[c.RoomId], c.Username)
		index := 0
		for _, v := range rooms[c.RoomId] {
			if v.Username != c.Username {
				rooms[c.RoomId][index] = v
				index++
			}
		}
		rooms[c.RoomId] = rooms[c.RoomId][:index]
	}
	mutex.Unlock()
}
func formatMessage(msg Message) []byte {
	data := make(map[string]interface{})
	switch msg.MsgType {
	case msgTypeOnline:
		data["username"] = msg.Data.(map[string]interface{})["username"]
		data["roomId"] = msg.Data.(map[string]interface{})["roomId"]
	case msgTypeSend:
		strUserId := msg.Data.(map[string]interface{})["userId"].(string)
		userId, _ := strconv.ParseUint(strUserId, 10, 32)
		username := model.FindUserById(uint(userId)).Username
		data["username"] = username
		data["roomId"] = msg.Data.(map[string]interface{})["roomId"]
		data["content"] = msg.Data.(map[string]interface{})["content"]
	}
	bytes, _ := json.Marshal(Message{
		MsgType: msg.MsgType,
		Data:    data,
	})
	return bytes
}
