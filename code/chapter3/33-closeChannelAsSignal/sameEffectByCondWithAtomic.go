package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var wg sync.WaitGroup
	c := sync.NewCond(&sync.Mutex{})
	var count int32
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			atomic.AddInt32(&count, 1)
			c.L.Lock()
			c.Wait()
			c.L.Unlock()
			fmt.Printf("%d has begin\n", i)
		}(i)
	}

	for {
		if atomic.LoadInt32(&count) == int32(5) {
			break
		}
	}

	fmt.Printf("unblocking goroutines...\n")
	c.Broadcast()
	wg.Wait()
}
