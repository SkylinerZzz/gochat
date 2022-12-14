package service

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gochat/util"
	"testing"
)

func TestBroadcaster(t *testing.T) {
	util.Init("../../config")
	rd := util.RedisPool.Get()
	defer rd.Close()
	list, err := redis.IntMap(rd.Do("hgetall", getUserListKey("1")))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(list["1"])
}
