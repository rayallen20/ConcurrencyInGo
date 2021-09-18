package main

import (
	"fmt"
	"sync"
	"time"
)

type data struct {
	mu sync.Mutex
	value int
}

func main() {
	var wg sync.WaitGroup

	printSum := func(d1, d2 *data) {
		defer wg.Done()

		d1.mu.Lock()
		defer d1.mu.Unlock()

		time.Sleep(2 * time.Second)

		d2.mu.Lock()
		defer d2.mu.Unlock()

		fmt.Printf("sum=%v\n", d1.value + d2.value)
	}

	var a, b data
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}
