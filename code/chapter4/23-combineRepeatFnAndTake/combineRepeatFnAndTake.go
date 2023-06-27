package main

import (
	"fmt"
	"math/rand"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	for value := range take(done, repeatFn(done, randFn), 10) {
		fmt.Println(value)
	}
}

func randFn() interface{} {
	return rand.Int()
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()

	return valueStream
}

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}
