package main

import "fmt"

func main() {
	fmt.Printf("fib(4) = %d\n", <-fib(4))
}

func fib(n int) <-chan int {
	result := make(chan int)

	go func() {
		defer close(result)

		if n <= 2 {
			result <- 1
			return
		}

		result <- <-fib(n-1) + <-fib(n-2)
	}()

	return result
}
