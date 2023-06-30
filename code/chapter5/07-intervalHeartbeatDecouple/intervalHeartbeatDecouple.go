package main

import "time"

func main() {

}

func doWork(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)

		for {
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
