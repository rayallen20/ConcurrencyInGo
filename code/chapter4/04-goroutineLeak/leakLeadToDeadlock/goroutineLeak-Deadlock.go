package main

import (
	"fmt"
	"sync"
)

func main() {
	doWork(nil)
	// 一些其他要执行的操作
	fmt.Println("done")
}

func doWork(stringStream <-chan string) <-chan interface{} {
	completed := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("doWork exited")
		defer close(completed)
		for s := range stringStream {
			// 模拟一些操作
			fmt.Printf("Reveived: %s\n", s)
		}
	}()
	wg.Wait()

	return completed
}
