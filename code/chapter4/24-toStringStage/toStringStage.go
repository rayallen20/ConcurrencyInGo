package main

func main() {
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
