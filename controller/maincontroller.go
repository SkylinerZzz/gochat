package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gochat/controller/session"
	"gochat/modelv2"
	"net/http"
	"regexp"
)

// LoginPage is responsible for displaying login page
func LoginPage(c *gin.Context) {
	// user has logged in recently
	if session.GetSession(c) != nil {
		c.Redirect(http.StatusFound, "/index")
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

// Login validates user info
func Login(c *gin.Context) {
	var u modelv2.UserInfo
	if err := c.ShouldBind(&u); err != nil {
		log.Errorf("[Login] failed to parse form, err = %s", err)
		c.JSON(http.StatusServiceUnavailable, nil)
		return
	}
	// input can not be empty
	log.WithFields(log.Fields{
		"username": u.Username,
		"password": u.Password,
	}).Info("[Login] someone try to login")

	if u.Username == "" {
		log.Info("[Login] user name is empty")
		//c.Writer.Write([]byte("<script>alert('please input username')</script>"))
		//c.HTML(http.StatusOK, "login.html", nil)
		c.JSON(http.StatusOK, gin.H{"msg": "Username is empty"})
		return
	}
	if u.Password == "" {
		log.Info("[Login] password is empty")
		//c.Writer.Write([]byte("<script>alert('please input password')</script>"))
		//c.HTML(http.StatusOK, "login.html", nil)
		c.JSON(http.StatusOK, gin.H{"msg": "Password is empty"})
		return
	}

	user, err := modelv2.FindUserByName(u.Username)
	if err != nil {
		log.Errorf("[Login] failed to find user by name, err = %s", err)
		c.JSON(http.StatusServiceUnavailable, nil)
		return
	}
	if user.ID > 0 && user.Password == u.Password {
		// login succeeded
		// store session
		session.SetSession(c, map[string]interface{}{"user_id": user.ID, "username": user.Username})
		//c.Redirect(http.StatusFound, "/index")
		c.JSON(http.StatusOK, gin.H{"msg": "success"})
		return
	} else {
		// login failed
		//c.Writer.Write([]byte("<script>alert('incorrect username or password')</script>"))
		//c.HTML(http.StatusOK, "login.html", nil)
		c.JSON(http.StatusOK, gin.H{"msg": "Incorrect username or password"})
		return
	}
}

// SignupPage is responsible for displaying login page
func SignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

// Signup validates user info format
func Signup(c *gin.Context) {
	var u modelv2.UserInfo
	if err := c.ShouldBind(&u); err != nil {
		log.Errorf("[Register] failed to parse form, err = %s", err)
		c.JSON(http.StatusServiceUnavailable, nil)
		return
	}
	log.WithFields(log.Fields{
		"username": u.Username,
		"password": u.Password,
	}).Info("someone trying to sign up")

	// validate input
	namePattern := "^[a-zA-Z0-9_-]{4,20}$"
	pwdPattern := "^[a-zA-Z0-9]{6,20}$"
	if m, _ := regexp.MatchString(namePattern, u.Username); !m {
		//c.Writer.Write([]byte("<script>alert('invalid username')</script>"))
		//c.HTML(http.StatusOK, "signup.html", nil)
		c.JSON(http.StatusOK, gin.H{"msg": "Invalid username"})
		return
	}
	if m, _ := regexp.MatchString(pwdPattern, u.Password); !m {
		//c.Writer.Write([]byte("<script>alert('invalid password')</script>"))
		//c.HTML(http.StatusOK, "signup.html", nil)
		c.JSON(http.StatusOK, gin.H{"msg": "Invalid password"})
		return
	}

	ex, err := modelv2.CheckUserExists(u.Username)
	if err != nil {
		log.Errorf("[Signup] failed to check whether username exists, err = %s", err)
		c.JSON(http.StatusServiceUnavailable, nil)
		return
	}
	if ex {
		// user exists
		log.Info("username already exists")
		//c.Writer.Write([]byte("<script>alert('username already exists')</script>"))
		//c.HTML(http.StatusOK, "signup.html", nil)
		c.JSON(http.StatusOK, gin.H{"msg": "Username already exists"})
		return
	} else {
		// sign up succeeded
		err = modelv2.AddUser(u)
		if err != nil {
			log.Errorf("[Signup] failed to add user, err = %s", err)
			c.JSON(http.StatusServiceUnavailable, nil)
			return
		}
		//c.Redirect(http.StatusFound, "/login")
		c.JSON(http.StatusOK, gin.H{"msg": "success"})
		return
	}
}

// Logout clears previous userinfo session
func Logout(c *gin.Context) {
	session.DelSession(c)
	c.Redirect(http.StatusFound, "/login")
}

// IndexPage is responsible for displaying index page
func IndexPage(c *gin.Context) {
	// get user info from session
	user := session.GetSession(c)
	if user == nil || (user["user_id"] == "" || user["username"] == "") {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	log.WithFields(log.Fields{
		"userId":   user["user_id"],
		"username": user["username"],
	}).Info("welcome to home")

	c.HTML(http.StatusOK, "index.html", gin.H{"username": user["username"]})
}

// NewRoomPage is responsible for displaying new room page
func NewRoomPage(c *gin.Context) {
	user := session.GetSession(c)
	if user == nil || (user["user_id"] == "" || user["username"] == "") {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "new.html", gin.H{
		"userId":   user["user_id"],
		"username": user["username"],
	})
}

// NewRoom creates new room
func NewRoom(c *gin.Context) {
	user := session.GetSession(c)
	var r modelv2.RoomInfo
	if err := c.ShouldBind(&r); err != nil {
		log.Errorf("[NewRoom] failed to parse form, err = %s", err)
		c.JSON(http.StatusServiceUnavailable, nil)
		return
	}
	log.WithFields(log.Fields{
		"userId":   user["user_id"],
		"username": user["username"],
	}).Info("someone try to create a room")

	namePattern := "^[a-zA-Z0-9_-]{1,20}$"
	if m, _ := regexp.MatchString(namePattern, r.RoomName); !m {
		//c.Writer.Write([]byte("<script>alert('invalid name')</script>"))
		//c.HTML(http.StatusOK, "new.html", nil)
		c.JSON(http.StatusOK, gin.H{"msg": "Invalid room name"})
		return
	}
	err := modelv2.AddRoom(r)
	if err != nil {
		log.Errorf("[NewRoom] failed to add room, err = %s", err)
		c.JSON(http.StatusServiceUnavailable, nil)
		return
	}

	log.Info(user["username"], " create a new room")
	//c.Redirect(http.StatusFound, "/index")
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

// Search lists rooms by name
func Search(c *gin.Context) {
	roomName := c.PostForm("room_name")
	if roomName == "" {
		c.JSON(http.StatusOK, gin.H{"msg": "Room name is empty"})
		return
	}
	rooms, err := modelv2.FindRoomsByName(roomName)
	if err != nil {
		log.Errorf("[Search] failed to find rooms by name, err = %s", err)
		c.JSON(http.StatusServiceUnavailable, nil)
		return
	}

	c.JSON(http.StatusOK, rooms)
}

// RoomPage is responsible for displaying room page
func RoomPage(c *gin.Context) {
	// store room name
	roomName := c.Query("room_name")
	session.UpdateSession(c, "room_name", roomName)

	user := session.GetSession(c)
	log.WithFields(log.Fields{
		"userId":   user["user_id"],
		"username": user["username"],
		"roomName": user["room_name"],
	}).Info("welcome to the room")

	c.HTML(http.StatusOK, "room.html", gin.H{
		"userId":    user["user_id"],
		"username":  user["username"],
		"roomId":    c.Param("roomId"),
		"room_name": user["room_name"],
	})
}

// PrivateChat gets or generates room id by users' id
func PrivateChat(c *gin.Context) {
	userId := c.PostForm("user_id")
	toUserId := c.PostForm("to_user_id")
	toUsername := c.PostForm("to_username")
	log.WithFields(log.Fields{
		"userId":     userId,
		"toUserId":   toUserId,
		"toUsername": toUsername,
	}).Info("someone try to chat with others")

	roomId, err := modelv2.GetPrivateRoomId(userId, toUserId)
	if err != nil {
		log.Errorf("[PrivateChat] failed to get private chat room id, err = %s", err)
		c.JSON(http.StatusServiceUnavailable, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"room_id": roomId})
}
