package main

import "fmt"

func main() {
	// chanOwner作为管道的拥有者 最终通过暴露一个只读管道
	// 将它的初始化和写入结果暴露给外部
	chanOwner := func() <-chan int {
		// 缓冲区长度为5的管道
		resultStream := make(chan int, 5)

		go func() {
			defer close(resultStream)
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()

		// 此处将一个双向管道隐式转换为了一个单向管道
		// 确保下游的消费者不会对该管道具有写权限
		return resultStream
	}

	// 不拥有管道的goroutine只需从管道中读取即可
	resultStream := chanOwner()
	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}

	fmt.Printf("Done receiving!\n")
}
