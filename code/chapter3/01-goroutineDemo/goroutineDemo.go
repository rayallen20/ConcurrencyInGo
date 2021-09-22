package main

import "fmt"

func main() {
	go sayHello()
	// 继续执行自己的逻辑
}

func sayHello() {
	fmt.Println("Hello")
}