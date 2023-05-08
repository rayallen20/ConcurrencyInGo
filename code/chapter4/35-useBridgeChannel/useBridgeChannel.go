package main

import "fmt"

func main() {
	for v := range bridge(nil, genValues()) {
		fmt.Printf("%v\n", v)
	}
}

func genValues() <-chan <-chan interface{} {
	chanStream := make(chan (<-chan interface{}))

	go func() {
		defer close(chanStream)

		for i := 0; i < 10; i++ {
			stream := make(chan interface{}, 1)
			stream <- i
			close(stream)
			chanStream <- stream
		}
	}()

	return chanStream
}

func bridge(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			var stream <-chan interface{}

			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}

				stream = maybeStream
			case <-done:
				return
			}

			for value := range orDone(done, stream) {
				select {
				case valueStream <- value:
				case <-done:
				}
			}
		}
	}()

	return valueStream
}

func orDone(done <-chan interface{}, c <-chan interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}

				select {
				case valueStream <- v:
				case <-done:
				}
			}
		}
	}()

	return valueStream
}
