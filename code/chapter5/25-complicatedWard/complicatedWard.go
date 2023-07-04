package main

import (
	"log"
	"time"
)

type startGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{})

func main() {
}

func doWorkFn(done <-chan interface{}, intList ...int) (startGoroutineFn, <-chan interface{}) {
	intChanStream := make(chan (<-chan interface{}))
	resultIntStream := bridge(done, intChanStream)

	doWork := func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
		intStream := make(chan interface{})
		heartbeat := make(chan interface{})

		go func() {
			defer close(intStream)
			select {
			case intChanStream <- intStream:
			case <-done:
				return
			}

			pulse := time.Tick(pulseInterval)

			for {
			valueLoop:
				for _, intVal := range intList {
					if intVal < 0 {
						log.Printf("negative value: %v\n", intVal)
						return
					}

					for {
						select {
						case <-pulse:
							sendPulse(heartbeat)
						case intStream <- intVal:
							continue valueLoop
						case <-done:
							return
						}
					}
				}
			}
		}()

		return heartbeat
	}

	return doWork, resultIntStream
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

func sendPulse(heartbeat chan interface{}) {
	select {
	case heartbeat <- struct{}{}:
	default:
	}
}
