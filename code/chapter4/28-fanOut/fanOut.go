package main

import (
	"math/rand"
	"runtime"
)

func main() {
	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
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

	toInt := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
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

	primeFinder := func(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
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

	randFn := func() interface{} { return rand.Intn(50000000) }

	done := make(chan interface{})
	defer close(done)

	randIntStream := toInt(done, repeatFn(done, randFn))

	// fan-out the stage of primeFinder()
	numFinders := runtime.NumCPU()
	finders := make([]<-chan interface{}, numFinders)
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}
}
