package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gochat/pkg/service"
)

func WsCreate(c *gin.Context) {
	wsConn, err := (&websocket.Upgrader{}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("[WsCreate] failed to upgrade protocol, err = %s", err)
		return
	}
	entry := service.NewEntry()
	roomId := c.Query("room_id")
	userId := c.Query("user_id")
	log.WithFields(log.Fields{
		"roomId": roomId,
		"userId": userId,
	}).Info("[WsCreate] create websocket connection")
	entry.Exec(wsConn, roomId, userId)
}
