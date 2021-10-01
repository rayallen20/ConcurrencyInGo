# sync.Cond

[原文出处](https://studygolang.com/articles/28072?fr=sidebar)

## 1. 使用场景

我需要完成一项任务,但是这项任务**需要满足一定条件才可以执行**,否则我就等待,直到满足条件为止.

那么我该如何获取到这个条件呢?

1. 循环获取
2. 当条件满足时,通知我

显然第2种方式效率高很多.

通知的方式,在GO中可以使用channel来实现.

```go
package main

import (
	"fmt"
	"time"
)

var mail = make(chan string)

func main() {
	go func() {
		<- mail
		fmt.Println("get chance to do something")
	}()

	time.Sleep(5 * time.Second)
	mail <- "this is a chance to do something"
	time.Sleep(2 * time.Second)
}
```

但是channel比较适用于1对1通知的方式,1对多通知并不是很适合.下面就来介绍另一种方式:`sync.Cond`

**`sync.Cond`就是用于实现条件变量的**,是基于`sync.Mutex`的基础上,增加了一个通知队列,通知的线程会从通知队列中唤醒1个或多个被通知的线程.

`sync.NewCond(&mutex)`:生成一个`cond`,需要传入一个`mutex`,因为阻塞等待通知的操作以及通知解除阻塞的操作就是基于`sync.Mutex`来实现的

`sync.Wait()`:用于等待通知

```go
func (c *Cond) Wait() {
	c.checker.check()
	t := runtime_notifyListAdd(&c.notify)
	c.L.Unlock()
	runtime_notifyListWait(&c.notify, t)
	c.L.Lock()
}
```

`Wait`自行解锁`c.L`并阻塞当前线程,在之后线程恢复执行时,`Wait`方法会在返回前锁定`c.L`.和其他系统不同,`Wait`除非被`Broatcast`或`Signal`唤醒,否则不会主动返回.

从命名就能看出,所有被`Wait`方法阻塞的线程,都被加入到了一个`notifyList`中.

`sync.Signal()`:用于发送单播

```go
func (c *Cond) Signal() {
	c.checker.check()
	runtime_notifyListNotifyOne(&c.notify)
}
```

`Signal`唤醒等待c的一个线程(如果存在).调用者在调用本方法时,建议(但并非必须)保持`c.L`的锁定

`runtime_notifyListNotifyOne(&c.notify)`:随机挑选一个协程进行通知,该线程`Wait`阻塞解除

`sync.Broatcast()`:用于广播

```go
func (c *Cond) Broadcast() {
	c.checker.check()
	runtime_notifyListNotifyAll(&c.notify)
}
```

`Broadcast`唤醒所有等待c的线程.调用者在调用本方法时,建议(但并非必须)保持`c.L`的锁定

`runtime_notifyListNotifyAll(&c.notify)`:通知所有等待的协程

## 2. 基本用法

例:主线程对多个goroutine的通知

```go
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

```

运行结果:

```
go run condDemo2.go
Signal 1...
3
Signal 2...
0
Broadcast...
9
4
8
1
5
7
2
6
```

## 3. 具体需求

现有4个worker和1个master,worker等待master分配指令,master持续在计数,计数到5时通知第1个worker,计数到10时通知第2和第3个worker.

有如下几种解决方案:

1. 所有worker循环去查看master的计数值,计数值满足自己的条件时,触发操作.
	- 缺点:无谓的消耗资源.因为worker需要一直占用CPU资源去查看master的计数值
2. 使用channel.有几个worker就有几个channel.worker1的协程里使用`<- channel`进行阻塞,当计数值到5时,master给worker1的channel发送信号.
	- 缺点:channel比较适用于1对1的场景.1对多时,需要创建很多channel,不是很美观.
3. 使用条件变量`sync.Cond`.针对多个worker,使用broadcast.由每个worker自行判断是否满足工作条件.

```go
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
```

运行结果:

```
go run condDemo3.go
worker1 started to work
worker1 work end
worker3 started to work
worker2 started to work
worker2 work end
worker3 work end
```

















































