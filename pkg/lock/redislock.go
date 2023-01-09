package lock

import (
	"github.com/gomodule/redigo/redis"
	"gochat/util"
	"time"
)

// Lock and Unlock are methods of distributed lock based on redis
func Lock(key string, timeout time.Duration) (bool, error) {
	rd := util.RedisPool.Get()
	defer rd.Close()

	desc, err := redis.String(rd.Do("set", key, time.Now(), "px", timeout.Milliseconds(), "nx"))
	if err != nil {
		return false, err
	}
	// get lock
	if desc == "OK" {
		return true, nil
	}
	// can not get lock
	return false, nil
}

func Unlock(key string) error {
	rd := util.RedisPool.Get()
	defer rd.Close()

	_, err := rd.Do("del", key)
	return err
}

// SpinLock wait for the lock
func SpinLock(key string, timeout time.Duration) (bool, error) {
	maxTry := 10
	for curTry := 0; curTry < maxTry; curTry++ {
		ok, err := Lock(key, timeout)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false, nil
}
