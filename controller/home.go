package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gochat/model"
	"gochat/session"
	"net/http"
	"regexp"
)

type Room struct {
	RoomName string `form:"roomName"`
}

func Index(c *gin.Context) {
	// get user info from session
	user := session.GetSession(c)
	logrus.WithFields(logrus.Fields{
		"userId":   user["userId"],
		"username": user["username"],
	}).Info("welcome to home")
	rs := model.ListAllRooms()
	c.HTML(http.StatusOK, "home.html", rs)
}
func NewPage(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", nil)
}
func New(c *gin.Context) {
	user := session.GetSession(c)
	var r Room
	if err := c.ShouldBind(&r); err != nil {
		logrus.Warn("binding failed: " + err.Error())
		return
	}
	logrus.WithFields(logrus.Fields{
		"userId":   user["userId"],
		"username": user["username"],
	}).Info("someone try to create a room")
	namePattern := "^[a-zA-Z0-9_-]{1,20}$"
	if m, _ := regexp.MatchString(namePattern, r.RoomName); !m {
		c.Writer.Write([]byte("<script>alert('invalid name')</script>"))
		c.HTML(http.StatusOK, "new.html", nil)
		return
	}
	model.AddRoom(map[string]interface{}{
		"userId":   user["userId"],
		"roomName": r.RoomName,
	})
	logrus.Info(user["username"], " create a new room")
	c.Redirect(http.StatusMovedPermanently, "/home")
}
