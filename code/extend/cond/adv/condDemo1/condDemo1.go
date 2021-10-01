package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

type FIFO struct {
	lock sync.Mutex
	cond *sync.Cond
	queue []int
}

type Queue interface {
	Pop() int
	Offer(num int) error
}

func (f *FIFO) Offer(num int) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.queue = append(f.queue, num)
	// 每向队列中添加1个元素 向所有Pop的goroutine广播1次
	f.cond.Broadcast()
	return nil
}

func (f *FIFO) Pop() int {
	f.lock.Lock()
	defer f.lock.Unlock()

	for len(f.queue) == 0 {
		f.cond.Wait()
	}

	item := f.queue[0]
	f.queue = f.queue[1:]
	return item
}

func main() {
	l := sync.Mutex{}
	fifo := &FIFO{
		lock: l,
		cond: sync.NewCond(&l),
		queue: []int{},
	}

	// 持续向队列投放数据
	go func() {
		for {
			fifo.Offer(rand.Int())
		}
	}()

	time.Sleep(1 * time.Second)

	// 持续从队列拿取数据
	go func() {
		for {
			fmt.Println(fmt.Sprintf("goroutine1 pop --> %d", fifo.Pop()))
		}
	}()

	// 持续从队列拿取数据
	go func() {
		for {
			fmt.Println(fmt.Sprintf("goroutine2 pop --> %d", fifo.Pop()))
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
