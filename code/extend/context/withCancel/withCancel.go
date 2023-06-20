package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	wg.Add(1)
	go doSelect(ctx)
	time.Sleep(3 * time.Second)
	cancel()
	wg.Wait()
}

func doSelect(ctx context.Context) {
LOOP:
	for {
		fmt.Printf("select data from DB\n")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LOOP
		default:
			continue
		}
	}
	wg.Done()
}
