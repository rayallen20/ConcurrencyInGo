package main

import "sync"

func main() {
	fanIn := func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})

		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		// 从所有channel中取值
		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		// 等待所有数据汇总完毕
		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}
}
