package main

import "fmt"

func main() {
	generator := func(done chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)

		go func() {
			defer close(intStream)
			for _, integer := range integers {
				select {
				case <-done:
					return
				case intStream <- integer:
				}
			}
		}()

		return intStream
	}

	multiply1 := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)

		go func() {
			defer close(multipliedStream)
			defer fmt.Printf("close multipliedStream 1\n")

			for integer := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- integer * multiplier:
				}
			}
		}()

		return multipliedStream
	}

	add2 := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)

		go func() {
			defer close(addedStream)
			defer fmt.Printf("close addedStream 2\n")

			for integer := range intStream {
				select {
				case <-done:
					return
				case addedStream <- integer + additive:
				}
			}
		}()

		return addedStream
	}

	multiply3 := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)

		go func() {
			defer close(multipliedStream)
			defer fmt.Printf("close multipliedStream 3\n")

			for integer := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- integer * multiplier:
				}
			}
		}()

		return multipliedStream
	}

	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply3(done, add2(done, multiply1(done, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}
