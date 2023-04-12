package main

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
}
