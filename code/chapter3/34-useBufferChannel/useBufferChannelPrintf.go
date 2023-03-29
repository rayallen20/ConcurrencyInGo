package main

import (
	"fmt"
)

func main() {
	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Printf("Producer done.\n")
		for i := 0; i < 5; i++ {
			fmt.Printf("Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Printf("Received: %d\n", integer)
	}
}
