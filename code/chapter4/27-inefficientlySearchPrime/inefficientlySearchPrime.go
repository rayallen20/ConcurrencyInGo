package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, randFn))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
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

func toInt(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
	intStream := make(chan int)

	go func() {
		defer close(intStream)

		for value := range valueStream {
			select {
			case <-done:
				return
			case intStream <- value.(int):
			}
		}
	}()

	return intStream
}

func primeFinder(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
	primeStream := make(chan interface{})

	go func() {
		defer close(primeStream)

		for integer := range intStream {
			// determine whether integer is prime
			integer -= 1
			prime := true
			for divisor := integer - 1; divisor > 1; divisor-- {
				if integer%divisor == 0 {
					prime = false
					break
				}
			}

			if prime {
				select {
				case <-done:
					return
				case primeStream <- integer:
				}
			}
		}
	}()

	return primeStream
}

func randFn() interface{} {
	return rand.Intn(50000000)
}
