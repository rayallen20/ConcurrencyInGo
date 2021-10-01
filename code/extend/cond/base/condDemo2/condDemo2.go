package main

import (
	"fmt"
	"sync"
	"time"
)

var locker sync.Mutex
var cond = sync.NewCond(&locker)

func main() {
	for i := 0; i < 10; i++ {
		go func(x int) {
			// 获取锁
			cond.L.Lock()

			// 释放锁
			defer cond.L.Unlock()

			// 等待通知 通知到来前阻塞当前goroutine
			cond.Wait()

			// 通知到来时 cond.Wait()会结束阻塞 后边就是当前线程要做的事情 此处仅打印
			fmt.Println(x)
		}(i)
	}

	// 睡眠1s 等待所有goroutine进入阻塞状态
	time.Sleep(1 * time.Second)
	fmt.Println("Signal 1...")

	// 1s后下发一个通知给已经获取锁的goroutine
	cond.Signal()

	time.Sleep(1 * time.Second)
	fmt.Println("Signal 2...")

	// 1s后再下发一个通知给已经获取锁的goroutine
	cond.Signal()

	time.Sleep(1 * time.Second)
	fmt.Println("Broadcast...")

	// 1s后发送广播给所有获取锁的goroutine
	cond.Broadcast()

	time.Sleep(1 * time.Second)
}
