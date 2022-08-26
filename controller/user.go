package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gochat/model"
	"gochat/session"
	"net/http"
	"regexp"
)

type User struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBind(&u); err != nil {
		logrus.Warn("binding failed: " + err.Error())
		return
	}
	// input can not be empty
	logrus.WithFields(logrus.Fields{
		"username": u.Username,
		"password": u.Password,
	}).Info("someone try to login")
	if u.Username == "" {
		logrus.Info("username is empty")
		c.Writer.Write([]byte("<script>alert('please input username')</script>"))
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}
	user := model.FindUserByName(u.Username)
	if u.Password == "" {
		logrus.Info("password is empty")
		c.Writer.Write([]byte("<script>alert('please input password')</script>"))
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}
	if user.ID > 0 && user.Password == u.Password {
		// login succeeded
		// store session
		session.SaveSession(c, user.ID)
		c.Redirect(http.StatusMovedPermanently, "/home")
		return
	} else {
		// login failed
		c.Writer.Write([]byte("<script>alert('incorrect username or password')</script>"))
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}
}

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func Register(c *gin.Context) {
	var u User
	if err := c.ShouldBind(&u); err != nil {
		logrus.Warn("binding failed: " + err.Error())
		return
	}
	// validate input
	logrus.WithFields(logrus.Fields{
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
		logrus.Info("username already exists")
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
