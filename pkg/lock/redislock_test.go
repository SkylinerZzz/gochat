package lock

import (
	"fmt"
	"gochat/util"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	util.Init("../../config")
	flag, _ := Lock("test", 10*time.Second)
	if flag {
		fmt.Println("get lock")
	} else {
		fmt.Println("can not get lock")
	}
	fmt.Println("release lock")
	Unlock("test")
	flag, _ = Lock("test", 10*time.Second)
	if flag {
		fmt.Println("get lock")
	} else {
		fmt.Println("can not get lock")
	}
}
