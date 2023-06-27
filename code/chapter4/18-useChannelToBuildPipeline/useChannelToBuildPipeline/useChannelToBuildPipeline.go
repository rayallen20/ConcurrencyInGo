package main

import "fmt"

func main() {
	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}

func generator(done <-chan interface{}, integers ...int) <-chan int {
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

func multiply(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
	multipliedStream := make(chan int)

	go func() {
		defer close(multipliedStream)
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

func add(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
	addedStream := make(chan int)

	go func() {
		defer close(addedStream)
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
