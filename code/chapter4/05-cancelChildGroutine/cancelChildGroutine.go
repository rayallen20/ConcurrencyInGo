package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		// 1s后取消goroutine的操作
		time.Sleep(1 * time.Second)
		fmt.Println("canceling doWork goroutine...")
		close(done)
	}()

	// 在关闭前阻塞
	<-terminated
	fmt.Println("done")
}

func doWork(done <-chan interface{}, stringStream <-chan string) <-chan interface{} {
	terminated := make(chan interface{})
	go func() {
		defer fmt.Println("doWork exited")
		defer close(terminated)
		for {
			select {
			case s := <-stringStream:
				// 模拟一些操作
				fmt.Printf("Reveived: %s\n", s)
			case <-done:
				return
			}
		}
	}()

	return terminated
}
