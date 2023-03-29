package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.L.Lock()
			c.Wait()
			fmt.Printf("%d has begin\n", i)
			c.L.Unlock()
		}(i)
	}

	fmt.Printf("unblocking goroutines...\n")
	c.Broadcast()
	wg.Wait()
}
