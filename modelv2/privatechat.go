package modelv2

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gochat/common"
	"gochat/pkg/lock"
	"gochat/util"
	"time"
)

// GetPrivateRoomId get room id which belongs to given two users
func GetPrivateRoomId(userId, toUserId string) (string, error) {
	rd := util.RedisPool.Get()
	defer rd.Close()

	// acquire lock before getting room id
	lockKey := getPrivateLockKey(userId, toUserId)
	ok, err := lock.SpinLock(lockKey, 10*time.Second)
	if err != nil {
		log.Errorf("[GetPrivateRoomId] failed to acquire spin lock, err = %s", err)
		return "", err
	}
	if !ok {
		return "", lock.ErrLockTimeout
	}
	defer lock.Unlock(lockKey)

	// get private chat room id
	key1 := getPrivateRoomKey(userId, toUserId)
	key2 := getPrivateRoomKey(toUserId, userId)
	roomId, err := redis.String(rd.Do("get", key1))
	if err != nil && !errors.Is(err, redis.ErrNil) {
		log.WithFields(log.Fields{
			"redisKey": key1,
		}).Errorf("failed to get redis value, err = %s", err)
		return "", err
	}
	if err == nil {
		return roomId, nil
	}

	// search another key
	roomId, err = redis.String(rd.Do("get", key2))
	if err != nil && !errors.Is(err, redis.ErrNil) {
		log.WithFields(log.Fields{
			"redisKey": key2,
		}).Errorf("failed to get redis value, err = %s", err)
		return "", err
	}
	if err == nil {
		_, err = rd.Do("set", key1, roomId)
		return roomId, err
	}

	// generate new room id
	roomId = uuid.New().String()
	// set key1 and key2
	_, err = rd.Do("set", key1, roomId)
	if err != nil {
		return "", err
	}
	_, err = rd.Do("set", key2, roomId)
	if err != nil {
		return "", err
	}
	return roomId, nil
}

func getPrivateRoomKey(userId, toUserId string) string {
	return common.PREFIX_PRIVATE_ROOM + userId + "to" + toUserId
}

func getPrivateLockKey(userId, toUserId string) string {
	if userId > toUserId {
		return toUserId + "_" + userId
	} else {
		return userId + "_" + toUserId
	}
}
