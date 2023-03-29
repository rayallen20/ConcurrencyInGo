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
		defer close(intStream)
		intStream <- 1
		wg.Done()
	}()

	go func() {
		integer, ok := <-intStream
		if !ok {
			fmt.Printf("intStream has been closed\n")
		} else {
			fmt.Printf("receive %d from intStream\n", integer)
		}
		wg.Done()
	}()

	wg.Wait()
}
