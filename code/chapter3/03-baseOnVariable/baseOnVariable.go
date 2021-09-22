package main

import "fmt"

func main() {
	sayHello := func() {
		fmt.Println("Hello")
	}

	go sayHello()
	// 继续执行自己的逻辑
}
