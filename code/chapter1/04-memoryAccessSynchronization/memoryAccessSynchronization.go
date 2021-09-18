package main

import (
	"fmt"
	"sync"
)

func main() {
	var memoryAccess sync.Mutex

	var data int

	go func() {
		memoryAccess.Lock()
		data++
		memoryAccess.Unlock()
	}()

	memoryAccess.Lock()
	if data == 0 {
		fmt.Printf("the value is %v\n", data)
	} else {
		fmt.Printf("the value is %v\n", data)
	}
	memoryAccess.Unlock()
}
