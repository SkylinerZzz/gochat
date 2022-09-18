package model

import (
	"fmt"
	"testing"
)

func TestAddUserToCache(t *testing.T) {
	var u User
	u = FindUserById(1)
	AddUserToCache(u)
	u, ok := FindUserByNameFromCache("Skyliner")
	if ok {
		fmt.Println("hit cache", u.Password)
	} else {
		fmt.Println("failed", u.CreatedAt)
	}
}
