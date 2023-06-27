package main

import "testing"

func main() {

}

func BenchmarkGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range toString(done, take(done, repeat(done, "a"), b.N)) {
	}
}

func BenchmarkTyped(b *testing.B) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range takeString(done, repeatString(done, "a"), b.N) {
	}
}

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
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

func repeatString(done <-chan interface{}, values ...string) <-chan string {
	valueStream := make(chan string)

	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
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

func takeString(done <-chan interface{}, valueStream <-chan string, num int) <-chan string {
	takeStream := make(chan string)

	go func() {
		defer close(takeStream)
		for i := num; i > 0 || i == -1; {
			if i != -1 {
				i--
			}
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func toString(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
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
