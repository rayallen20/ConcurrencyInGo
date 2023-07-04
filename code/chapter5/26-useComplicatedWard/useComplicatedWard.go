package main

import (
	"log"
	"os"
	"time"
)

type startGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{})

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	done := make(chan interface{})
	defer close(done)

	doWork, resultIntStream := doWorkFn(done, 1, 2, -1, 3, 4, 5)
	doWorkWithSteward := steward(1*time.Second, doWork)
	doWorkWithSteward(done, 1*time.Hour)

	for intVal := range take(done, resultIntStream, 6) {
		log.Printf("Received: %v\n", intVal)
	}
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

func steward(timeout time.Duration, startGoroutine startGoroutineFn) startGoroutineFn {
	return func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
		heartbeat := make(chan interface{})

		go func() {
			defer close(heartbeat)

			wardDone, wardHeartbeat := startWard(done, startGoroutine, timeout)
			pulse := time.Tick(pulseInterval)

		monitorLoop:
			for {
				timeoutSignal := time.After(timeout)
				for {
					select {
					case <-pulse:
						sendPulse(heartbeat)
					case <-wardHeartbeat:
						continue monitorLoop
					case <-timeoutSignal:
						log.Println("steward: ward unhealthy; restarting")
						close(wardDone)
						wardDone, wardHeartbeat = startWard(done, startGoroutine, timeout)
						continue monitorLoop
					case <-done:
						return
					}
				}
			}
		}()

		return heartbeat
	}
}

func startWard(done <-chan interface{}, startGoroutine startGoroutineFn, timeout time.Duration) (chan interface{}, <-chan interface{}) {
	wardDone := make(chan interface{})
	wardHeartbeat := startGoroutine(or(wardDone, done), timeout/2)
	return wardDone, wardHeartbeat
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
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
