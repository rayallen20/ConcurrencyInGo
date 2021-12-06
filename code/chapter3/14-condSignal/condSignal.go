package main

import (
	"fmt"
	"sync"
	"time"
)

var i int = 0

func main() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	// 移除元素 用于模拟发射信号一方的goroutine
	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Printf("Removed from queue %d\n", i)
		c.L.Unlock()
		c.Signal()
	}

	// 推送元素 用于模拟等待并接收信号一方的goroutine
	for i = 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Printf("Adding to queue %d\n", i)
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
}
