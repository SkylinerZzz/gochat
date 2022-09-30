package util

import "sync"

type SafeMap interface {
	Read(interface{}) (interface{}, bool)
	Write(key, val interface{})
	Delete(interface{})
}

type UserMap struct {
	mu   sync.RWMutex
	data map[string]bool
}

func (u *UserMap) Read(key interface{}) (interface{}, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	val, ok := u.data[key.(string)]
	return val, ok
}
func (u *UserMap) Write(key, val interface{}) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.data[key.(string)] = val.(bool)
}
func (u *UserMap) Delete(key interface{}) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.data, key.(string))
}
func NewUserMap() *UserMap {
	return &UserMap{
		mu:   sync.RWMutex{},
		data: map[string]bool{},
	}
}
