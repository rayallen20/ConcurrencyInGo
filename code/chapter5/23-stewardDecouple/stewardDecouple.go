package main

import (
	"log"
	"time"
)

// startGoroutineFn 启动被监控的goroutine 在监控到该goroutine停止工作时重新启动一个goroutine并监控该goroutine
type startGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{})

func main() {
}

// steward 返回一个函数 该函数启动并监控goroutine 若该goroutine停止工作 则该函数重新启动一个goroutine并监控该goroutine
// timeout 被监控的goroutine的超时时间
// startGoroutine 启动被监控的goroutine的函数
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

func sendPulse(heartbeat chan interface{}) {
	select {
	case heartbeat <- struct{}{}:
	default:
	}
}
