package task

import "errors"

var (
	ErrUnknownWsMessageType = errors.New("unknown type of WsMessage")
)
