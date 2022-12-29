package router

import (
	"github.com/gin-gonic/gin"
	"gochat/controller"
	"gochat/controller/session"
)

func Init(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("view/*")
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/icon/favicon.ico")
	sr := r.Group("/", session.EnableSession())
	{
		sr.GET("/login", controller.LoginPage)
		sr.POST("/login", controller.Login)
		sr.GET("/register", controller.RegisterPage)
		sr.POST("/register", controller.Register)
		sr.GET("/home", controller.IndexPage)
		sr.GET("/home/new", controller.NewRoomPage) // create new room
		sr.POST("/home/new", controller.NewRoom)
		sr.GET("/room/:roomId", controller.RoomPage) // enter into the room
		sr.GET("/room/ws", controller.WsCreate)      // build websocket connection
	}
	return r
}
