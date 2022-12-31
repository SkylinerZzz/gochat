package session

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func EnableSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte("SkylinerZzz"))
	// register map[string]interface{}
	gob.Register(map[string]interface{}{})
	return sessions.Sessions("UserInfo", store)
}

func SetSession(c *gin.Context, info map[string]interface{}) {
	s := sessions.Default(c)
	s.Set("UserInfo", info)
	err := s.Save()
	if err != nil {
		log.Errorf("session: failed to save user info, err = %s", err)
	}
}

func GetSession(c *gin.Context) map[string]interface{} {
	s := sessions.Default(c)
	v := s.Get("UserInfo")
	if v == nil {
		return nil
	}
	info := v.(map[string]interface{})
	return info
}

func DelSession(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	err := s.Save()
	if err != nil {
		log.Errorf("session: failed to delete user info, err = %s", err)
	}
}
