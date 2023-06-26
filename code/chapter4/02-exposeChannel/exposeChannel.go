package main

import "fmt"

func main() {
	results := chanOwner()
	consumer(results)
}

func chanOwner() <-chan int {
	results := make(chan int, 5)
	go func() {
		defer close(results)
		for i := 0; i < 5; i++ {
			results <- i
		}
	}()
	return results
}

func consumer(results <-chan int) {
	for result := range results {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Printf("Done receiving\n")
}
