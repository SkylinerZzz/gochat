package router

import (
	"github.com/gin-gonic/gin"
	"gochat/controller"
	"gochat/controller/session"
)

func Init() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("view/template/*")
	r.Static("/view/static", "./view/static")
	r.Static("/view/bootstrap-5.3.0-alpha1-dist", "./view/bootstrap-5.3.0-alpha1-dist")
	r.StaticFile("/favicon.ico", "./view/static/icon/favicon.ico")

	sr := r.Group("/", session.EnableSession())
	{
		sr.GET("/login", controller.LoginPage)
		sr.POST("/login", controller.Login)
		sr.GET("/signup", controller.SignupPage)
		sr.POST("/signup", controller.Signup)
		sr.GET("/logout", controller.Logout)
		sr.GET("/index", controller.IndexPage)
		sr.GET("/index/new", controller.NewRoomPage) // create new room
		sr.POST("/index/new", controller.NewRoom)
		sr.GET("/room/:roomId", controller.RoomPage) // enter into the room
		sr.GET("/room/ws", controller.WsCreate)      // build websocket connection
	}
	return r
}
