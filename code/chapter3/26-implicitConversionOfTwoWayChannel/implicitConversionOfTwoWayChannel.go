package main

import "fmt"

func main() {
	var receiveChan <-chan interface{}
	var sendChan chan<- interface{}
	dataStream := make(chan interface{})

	receiveChan = dataStream // 隐式转换:将双向管道转换为只读管道
	sendChan = dataStream    // 隐式转换:将双向管道转换为只写管道

	fmt.Printf("%#v\n", receiveChan)
	fmt.Printf("%#v\n", sendChan)
}
