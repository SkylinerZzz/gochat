package common

import (
	"fmt"
	"sync"
	"testing"
)

func TestClientMap(t *testing.T) {
	ClientMap["room_1"] = &sync.Map{}
	ClientMap["room_1"].Store("user_1", "wy")
	fmt.Println(ClientMap["room_1"].Load("user_1"))
}
