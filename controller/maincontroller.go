package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gochat/controller/session"
	"gochat/model"
	"net/http"
	"regexp"
)

// UserInfo describes user info from registration and login
type UserInfo struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

// RoomInfo describes room info while creating new room
type RoomInfo struct {
	RoomName string `form:"room_name"`
}

// LoginPage is responsible for displaying login page
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Login validates user info
func Login(c *gin.Context) {
	var u UserInfo
	if err := c.ShouldBind(&u); err != nil {
		log.Errorf("[Login] failed to parse form, err = %s", err)
		return
	}
	// input can not be empty
	log.WithFields(log.Fields{
		"username": u.Username,
		"password": u.Password,
	}).Info("[Login] someone try to login")

	if u.Username == "" {
		log.Info("[Login] user name is empty")
		c.Writer.Write([]byte("<script>alert('please input username')</script>"))
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	user := model.FindUserByName(u.Username)
	if u.Password == "" {
		log.Info("password is empty")
		c.Writer.Write([]byte("<script>alert('please input password')</script>"))
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}
	if user.ID > 0 && user.Password == u.Password {
		// login succeeded
		// store session
		session.SetSession(c, map[string]interface{}{"user_id": user.ID, "username": user.Username})
		c.Redirect(http.StatusMovedPermanently, "/home")
		return
	} else {
		// login failed
		c.Writer.Write([]byte("<script>alert('incorrect username or password')</script>"))
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}
}

// RegisterPage is responsible for displaying login page
func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

// Register validates user info format
func Register(c *gin.Context) {
	var u UserInfo
	if err := c.ShouldBind(&u); err != nil {
		log.Errorf("[Register] failed to parse form, err = %s", err)
		return
	}
	// validate input
	log.WithFields(log.Fields{
		"username": u.Username,
		"password": u.Password,
	}).Info("someone trying to sign up")

	namePattern := "^[a-zA-Z0-9_-]{4,20}$"
	pwdPattern := "^[a-zA-Z0-9]{6,20}$"
	if m, _ := regexp.MatchString(namePattern, u.Username); !m {
		c.Writer.Write([]byte("<script>alert('invalid username')</script>"))
		c.HTML(http.StatusOK, "register.html", nil)
		return
	}
	if m, _ := regexp.MatchString(pwdPattern, u.Password); !m {
		c.Writer.Write([]byte("<script>alert('invalid password')</script>"))
		c.HTML(http.StatusOK, "register.html", nil)
		return
	}
	user := model.FindUserByName(u.Username)
	if user.ID > 0 {
		// user exists
		log.Info("username already exists")
		c.Writer.Write([]byte("<script>alert('username already exists')</script>"))
		c.HTML(http.StatusOK, "register.html", nil)
		return
	} else {
		// sign up succeeded
		model.AddUser(map[string]interface{}{
			"username": u.Username,
			"password": u.Password,
		})
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
}

// IndexPage is responsible for displaying index page
func IndexPage(c *gin.Context) {
	// get user info from session
	user := session.GetSession(c)
	log.WithFields(log.Fields{
		"userId":   user["user_id"],
		"username": user["username"],
	}).Info("welcome to home")

	rs := model.ListAllRooms()
	c.HTML(http.StatusOK, "home.html", rs)
}

// NewRoomPage is responsible for displaying new room page
func NewRoomPage(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", nil)
}

// NewRoom creates new room
func NewRoom(c *gin.Context) {
	user := session.GetSession(c)
	var r RoomInfo
	if err := c.ShouldBind(&r); err != nil {
		log.Errorf("[NewRoom] failed to parse form, err = %s", err)
		return
	}
	log.WithFields(log.Fields{
		"userId":   user["user_id"],
		"username": user["username"],
	}).Info("someone try to create a room")

	namePattern := "^[a-zA-Z0-9_-]{1,20}$"
	if m, _ := regexp.MatchString(namePattern, r.RoomName); !m {
		c.Writer.Write([]byte("<script>alert('invalid name')</script>"))
		c.HTML(http.StatusOK, "new.html", nil)
		return
	}
	model.AddRoom(map[string]interface{}{
		"userId":   user["user_id"],
		"roomName": r.RoomName,
	})
	log.Info(user["username"], " create a new room")
	c.Redirect(http.StatusMovedPermanently, "/home")
}

// RoomPage is responsible for displaying room page
func RoomPage(c *gin.Context) {
	user := session.GetSession(c)
	log.WithFields(log.Fields{
		"userId":   user["user_id"],
		"username": user["username"],
	}).Info("welcome to the room")

	c.HTML(http.StatusOK, "room.html", gin.H{
		"userId":   user["user_id"],
		"username": user["username"],
		"roomId":   c.Param("roomId"),
	})
}
