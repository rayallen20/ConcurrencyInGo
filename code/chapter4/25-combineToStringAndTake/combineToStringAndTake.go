package main

import "fmt"

func main() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})

		go func() {
			defer close(valueStream)

			for {
				for _, value := range values {
					select {
					case <-done:
						return
					case valueStream <- value:
					}
				}
			}
		}()

		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
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

	toString := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
		stringStream := make(chan string)

		go func() {
			defer close(stringStream)

			for value := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- value.(string):
				}
			}
		}()

		return stringStream
	}

	done := make(chan interface{})
	defer close(done)

	var message string

	for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		message += token
	}

	fmt.Printf("message = %s\n", message)
}
