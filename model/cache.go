package model

import (
	"encoding/json"
	"fmt"
)

func AddUserToCache(u User) {
	// set key by id
	key := fmt.Sprintf("user_id: %v", u.ID)
	val, _ := json.Marshal(u)
	ChatCache.Do("set", key, val, "ex", 3600)
	// set key by name
	key = fmt.Sprintf("user_name: %v", u.Username)
	ChatCache.Do("set", key, val, "ex", 3600)
}
func FindUserByIdFromCache(id uint) (User, bool) {
	var u User
	var ok bool
	key := fmt.Sprintf("user_id: %v", id)
	val, _ := ChatCache.Do("get", key)
	if val != nil {
		// hit cache
		json.Unmarshal(val.([]byte), &u)
		ok = true
	} else {
		ok = false
	}
	return u, ok
}
func FindUserByNameFromCache(name string) (User, bool) {
	var u User
	var ok bool
	key := fmt.Sprintf("user_name: %v", name)
	val, _ := ChatCache.Do("get", key)
	if val != nil {
		// hit cache
		json.Unmarshal(val.([]byte), &u)
		ok = true
	} else {
		ok = false
	}
	return u, ok
}
