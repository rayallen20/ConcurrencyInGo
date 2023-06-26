package main

import (
	"fmt"
	"math/rand"
)

func main() {
	randStream := newRandStream()
	fmt.Println("3 random integers:")

	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}

func newRandStream() <-chan int {
	randStream := make(chan int)
	go func() {
		defer fmt.Println("newRandStream closure exited.")
		defer close(randStream)
		for {
			randStream <- rand.Int()
		}
	}()

	return randStream
}
