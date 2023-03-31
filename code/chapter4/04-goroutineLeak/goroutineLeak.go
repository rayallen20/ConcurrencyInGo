package main

import (
	"fmt"
)

func main() {
	doWork := func(stringStream <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(completed)
			for s := range stringStream {
				// 模拟一些操作
				fmt.Printf("Reveived: %s\n", s)
			}
		}()

		return completed
	}

	doWork(nil)
	// 一些其他要执行的操作
	fmt.Println("done")
}
