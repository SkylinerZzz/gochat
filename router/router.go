package router

import (
	"github.com/gin-gonic/gin"
	"gochat/controller"
	"gochat/session"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("view/*")
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/icon/favicon.ico")
	sr := r.Group("/", session.EnableSession())
	{
		sr.GET("/login", controller.LoginPage)
		sr.POST("/login", controller.Login)
		sr.GET("/register", controller.RegisterPage)
		sr.POST("/register", controller.Register)
		sr.GET("/home", controller.Index)
		sr.GET("/home/new", controller.NewPage) // create new room
		sr.POST("/home/new", controller.New)
		sr.GET("/room/:roomId", controller.Enter) // enter into the room
		sr.GET("/room/ws", controller.Run)        // build websocket connection
	}
	return r
}
