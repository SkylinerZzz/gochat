package service

import (
	"encoding/json"
	"fmt"
	"gochat/common"
	"gochat/modelv2"
	"gochat/util"
	"testing"
)

func TestSub(t *testing.T) {
	util.Init("../../config")
	data := modelv2.Message{}
	wsMessage := common.WsMessage{Data: &data}
	raw := `
		{
			"type":1,
			"data":{
				"user_id":"123456",
				"room_id":"123456",
				"content":"hello world"
			}
		}
	`
	err := json.Unmarshal([]byte(raw), &wsMessage)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(wsMessage, data)
	sub := NewSubscriber()
	sub.Exec("1")
}
