package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}

func doWork(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func() {
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)

		for i := 0; i < 2; i++ {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse(heartbeat)
			case r := <-workGen:
				sendResult(done, results, r, pulse, heartbeat)
			}
		}
	}()

	return heartbeat, results
}

func sendPulse(heartbeat chan interface{}) {
	select {
	case heartbeat <- struct{}{}:
	default:
	}
}

func sendResult(done <-chan interface{}, results chan time.Time, r time.Time, pulse <-chan time.Time, heartbeat chan interface{}) {
	for {
		select {
		case <-done:
			return
		case <-pulse:
			sendPulse(heartbeat)
		case results <- r:
			return
		}
	}
}
