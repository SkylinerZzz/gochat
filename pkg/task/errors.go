package task

import "errors"

var (
	ErrUnknownWsMessageType = errors.New("queuetask: unknown type of WsMessage")
)
