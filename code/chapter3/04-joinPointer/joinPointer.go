package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	sayHello := func() {
		defer wg.Done()
		fmt.Println("Hello")
	}

	wg.Add(1)
	go sayHello()
	// 继续执行自己的逻辑
	wg.Wait()
}
