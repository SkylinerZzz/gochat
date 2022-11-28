package adapter

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, _ := context.WithTimeout(ctx, 5*time.Second)
	go doContext(c)
	time.Sleep(3 * time.Second)
}

func doContext(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			return
		default:
			fmt.Println("run")
			time.Sleep(1 * time.Second)
		}
	}
}
