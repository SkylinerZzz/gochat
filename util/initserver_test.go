package util

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	loadConfig("../config")
	fmt.Println(Config)
}

func TestInit(t *testing.T) {
	loadConfig("../config")
	initRedis()
	initDB()
	conn := RedisPool.Get()
	_, err := conn.Do("set", "name", "Skyliner", "ex", 5)
	if err != nil {
		t.Fatal(err)
	}
}
