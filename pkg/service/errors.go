package service

import "errors"

var (
	ErrInvalidParams = errors.New("service: invalid parameters")
	ErrWsConnClosed  = errors.New("websocket: close 1005 (no status)")
)
