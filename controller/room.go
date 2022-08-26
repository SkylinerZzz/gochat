package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gochat/session"
	"net/http"
)

func Enter(c *gin.Context) {
	user := session.GetSession(c)
	logrus.WithFields(logrus.Fields{
		"userId":   user["userId"],
		"username": user["username"],
	}).Info("welcome to the room")
	c.HTML(http.StatusOK, "room.html", gin.H{
		"userId":   user["userId"],
		"username": user["username"],
		"roomId":   c.Param("roomId"),
	})
}
