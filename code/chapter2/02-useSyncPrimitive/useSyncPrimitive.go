package main

import (
	"fmt"
	"sync"
)

type Data struct {
	mu sync.Mutex
	data int
}

func (d *Data) Increment() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data++
}

func main() {
	var data Data = Data{
		mu: sync.Mutex{},
		data: 0,
	}

	go func() {
		data.Increment()
	}()

	if data.data == 0 {
		fmt.Printf("the value is %v\n", data.data)
	} else {
		fmt.Printf("the value is %v\n", data.data)
	}
}
