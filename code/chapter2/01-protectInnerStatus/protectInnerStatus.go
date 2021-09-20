package main

import "sync"

type Counter struct {
	mu sync.Mutex
	data int
}

func (c *Counter) Increment()  {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data++
}

func main() {

}
