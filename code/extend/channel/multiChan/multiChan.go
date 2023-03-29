package main

import (
	"fmt"
	"sync"
)

func main() {
	var intStream chan int
	intStream = make(chan int)
	wg := &sync.WaitGroup{}
	producer := func() {
		for i := 0; i < 10; i++ {
			intStream <- i
		}
		close(intStream)
		wg.Done()
	}

	consumer := func(declare string) {
		for value := range intStream {
			fmt.Printf("%s get value from chan = %d\n", declare, value)
		}
		wg.Done()
	}

	wg.Add(3)
	go producer()
	go consumer("consumer1")
	go consumer("consumer2")
	wg.Wait()
}
