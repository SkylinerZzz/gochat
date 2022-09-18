package model

import "github.com/gomodule/redigo/redis"

var ChatCache redis.Conn

func init() {
	ChatCache, _ = redis.Dial("tcp", "localhost:6379", redis.DialPassword("123456"))
}
