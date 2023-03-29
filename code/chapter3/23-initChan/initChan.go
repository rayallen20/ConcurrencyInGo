package main

import "fmt"

func main() {
	var dataStream chan interface{}     // 声明一个channel.由于声明的类型为空接口,因此通常说它的类型是interface{}
	dataStream = make(chan interface{}) // 使用内置函数make()实例化channel
	fmt.Printf("%v\n", dataStream)
}
