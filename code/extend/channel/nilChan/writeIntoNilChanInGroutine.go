package main

import (
	"fmt"
	"sync"
)

func main() {
	var intStream chan int
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		intStream <- 1
		close(intStream)
	}()

	go func() {
		for integer := range intStream {
			fmt.Printf("receive %d from intStream\n", integer)
		}
		wg.Done()
	}()

	wg.Wait()
}
