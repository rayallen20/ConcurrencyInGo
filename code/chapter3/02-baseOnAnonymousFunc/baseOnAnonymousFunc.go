package main

import "fmt"

func main() {
	go func() {
		fmt.Println("Hello")
	}()
	// 继续执行自己的逻辑
}
