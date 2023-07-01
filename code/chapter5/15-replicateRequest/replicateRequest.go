package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	done := make(chan interface{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go doWork(done, i, &wg, result)
	}

	firstReturned := <-result
	close(done)
	wg.Wait()

	fmt.Printf("Received an answer from #%v\n", firstReturned)
}

func doWork(done <-chan interface{}, id int, wg *sync.WaitGroup, result chan<- int) {
	started := time.Now()
	defer wg.Done()

	// 模拟随机时长的负载
	simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second

	select {
	case <-done:
	case <-time.After(simulatedLoadTime):
	}

	select {
	case <-done:
	case result <- id:
	}

	took := time.Since(started)
	// 显示处理程序假定需要的时长(因为受done channel的控制 不一定真的会走第36行的time.After分支)
	if took < simulatedLoadTime {
		took = simulatedLoadTime
	}
	fmt.Printf("%v took %v\n", id, took)
}
