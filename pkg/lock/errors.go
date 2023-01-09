package lock

import "errors"

var (
	ErrLockTimeout = errors.New("lock: lock acquisition timeout")
)
