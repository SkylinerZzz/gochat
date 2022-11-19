package util

import (
	"fmt"
	"testing"
)

func Test_loadConfig(t *testing.T) {
	loadConfig("../config")
	fmt.Println(Config)
}

func Test_init(t *testing.T) {
	loadConfig("../config")
	initRedis()
	initDB()
	conn := RedisPool.Get()
	_, err := conn.Do("set", "name", "Skyliner", "ex", 5)
	if err != nil {
		t.Fatal(err)
	}
}
