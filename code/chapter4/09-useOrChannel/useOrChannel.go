package main

import (
	"fmt"
	"time"
)

func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		// 没有信号管道 直接返回
		case 0:
			return nil
		// 只有1个信号管道 返回该信号管道
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		// 有2个或以上的信号管道
		go func() {
			defer close(orDone)

			switch len(channels) {
			// 有2个信号管道
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}

			// 有2个以上的信号管道
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

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v\n", time.Since(start))
}
