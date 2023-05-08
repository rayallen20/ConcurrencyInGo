package main

func main() {

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
