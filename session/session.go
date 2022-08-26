package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gochat/model"
)

func EnableSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte("cookie_key"))
	return sessions.Sessions("userinfo", store)
}
func SaveSession(c *gin.Context, value interface{}) {
	session := sessions.Default(c)
	session.Set("userId", value)
	session.Save()
}
func GetSession(c *gin.Context) map[string]interface{} {
	session := sessions.Default(c)
	userId := session.Get("userId").(uint)
	data := make(map[string]interface{})
	u := model.FindUserById(userId)
	data["userId"] = u.ID
	data["username"] = u.Username
	return data
}
