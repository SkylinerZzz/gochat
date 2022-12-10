package lock

import "time"

// Lock and Unlock are methods of distributed lock based on redis
func Lock(key string, timeout time.Duration) (bool, error) {
	return false, nil
}

func Unlock(key string) error {
	return nil
}
