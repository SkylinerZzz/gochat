package model

import (
	"fmt"
	"testing"
)

func TestAddUserToCache(t *testing.T) {
	rs := ListAllRoomsFromCache()
	for _, r := range rs {
		fmt.Println(r.RoomName)
	}
}
