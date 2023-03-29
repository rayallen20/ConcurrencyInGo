package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	begin := make(chan interface{})
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%d has begin\n", i)
		}(i)
	}

	fmt.Printf("unblocking goroutines...\n")
	close(begin)
	wg.Wait()
}

//
