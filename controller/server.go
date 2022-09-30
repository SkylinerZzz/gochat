package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gochat/model"
	"gochat/util"
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

type ClientSlice struct { // safe client slice
	mu   sync.RWMutex
	data []Client
}

func NewRoomSlice() *ClientSlice {
	return &ClientSlice{
		mu:   sync.RWMutex{},
		data: []Client{},
	}
}
func (r *ClientSlice) Remove(username string) { // remove client mapping
	r.mu.Lock()
	defer r.mu.Unlock()
	index := 0
	for _, v := range r.data {
		if v.Username != username {
			r.data[index] = v
			index++
		}
	}
	r.data = r.data[:index]
}
func (r *ClientSlice) Append(c Client) { // add client mapping
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = append(r.data, c)
}
func (r *ClientSlice) Broadcast(msg Message) { // broadcast message in certain room
	clients := r.data
	for _, c := range clients {
		err := c.Conn.WriteMessage(websocket.TextMessage, formatMessage(msg))
		if err != nil {
			logrus.Warn(err)
		}
	}
}

var (
	once    = sync.RWMutex{}                 // ensure users[string] should be initialized once
	rooms   = make(map[string]*ClientSlice)  // map clients to the room
	users   = make(map[string]*util.UserMap) //user mapping, avoid duplicate connections
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
			if users[e.RoomId] != nil && rooms[e.RoomId] != nil {
				entering(e)
			}
			once.RUnlock()
			once.Lock()
			if users[e.RoomId] == nil {
				users[e.RoomId] = util.NewUserMap()
			}
			if rooms[e.RoomId] == nil {
				rooms[e.RoomId] = NewRoomSlice()
			}
			entering(e)
			once.Unlock()
		case l := <-leave:
			leaving(l)
		case msg := <-message:
			logrus.Info("broadcasting ...")
			roomId := msg.Data.(map[string]interface{})["roomId"].(string)
			if rooms[roomId] != nil {
				rooms[roomId].Broadcast(msg)
			} else {
				logrus.Warn("room map not be initialized yet!")
			}
		case <-done:
			return
		}
	}
}
func entering(c Client) {
	// update users mapping while entering
	if _, ok := users[c.RoomId].Read(c.Username); !ok {
		users[c.RoomId].Write(c.Username, true)
		rooms[c.RoomId].Append(c)
	}
	logrus.WithFields(logrus.Fields{
		"roomId":   c.RoomId,
		"roomSize": len(rooms[c.RoomId].data),
		"username": c.Username,
	}).Info("an user enter into the room")
}
func leaving(c Client) {
	// update users mapping while leaving
	if _, ok := users[c.RoomId].Read(c.Username); ok {
		users[c.RoomId].Delete(c.Username)
		rooms[c.RoomId].Remove(c.Username)
	}
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
