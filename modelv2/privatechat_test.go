package modelv2

import (
	"fmt"
	"testing"
)

func TestLockKey(t *testing.T) {
	fmt.Println(getPrivateLockKey("3", "2"))
}
