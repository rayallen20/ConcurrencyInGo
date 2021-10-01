package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	mutex := sync.Mutex{}
	cond := sync.NewCond(&mutex)
	mail := 1

	// master
	go func() {
		for count := 0; count <= 15; count++ {
			time.Sleep(1 * time.Second)
			mail = count
			cond.Broadcast()
		}
	}()

	// worker1
	go func() {
		// 触发条件:计数器 = 5
		for mail != 5 {
			cond.L.Lock()
			// 若计数器 != 5 则进入cond.Wait()等待
			// 当cond.Broadcast()通知传递过来时 wait阻塞解除
			// 进入下一次循环
			cond.Wait()
			cond.L.Unlock()
		}

		// 当mail == 5时 跳出循环 开始工作
		fmt.Println("worker1 started to work")
		time.Sleep(3 * time.Second)
		fmt.Println("worker1 work end")
	}()

	// worker2
	go func() {
		// 触发条件:计数器 = 10
		for mail != 10 {
			cond.L.Lock()
			// 若计数器 != 10 则进入cond.Wait()等待
			// 当cond.Broadcast()通知传递过来时 wait阻塞解除
			// 进入下一次循环
			cond.Wait()
			cond.L.Unlock()
		}

		// 当mail == 10时 跳出循环 开始工作
		fmt.Println("worker2 started to work")
		time.Sleep(3 * time.Second)
		fmt.Println("worker2 work end")
	}()

	// worker3
	go func() {
		// 触发条件:计数器 = 10
		for mail != 10 {
			cond.L.Lock()
			// 若计数器 != 10 则进入cond.Wait()等待
			// 当cond.Broadcast()通知传递过来时 wait阻塞解除
			// 进入下一次循环
			cond.Wait()
			cond.L.Unlock()
		}

		// 当mail == 10时 跳出循环 开始工作
		fmt.Println("worker3 started to work")
		time.Sleep(3 * time.Second)
		fmt.Println("worker3 work end")
	}()

	// worker4
	go func() {
		// 无论何时都不工作
	}()

	time.Sleep(20 * time.Second)
}
