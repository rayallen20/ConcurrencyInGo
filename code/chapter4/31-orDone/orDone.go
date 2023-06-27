package main

import "fmt"

func main() {
	done := make(chan interface{})
	defer close(done)
	myChan := make(chan interface{})

	for value := range orDone(done, myChan) {
		// 使用从channel中读取到的值执行一些逻辑
		fmt.Printf("%#v\n", value)
	}
}

func orDone(done <-chan interface{}, c <-chan interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				// 若给定的channel关闭 也直接返回
				if !ok {
					return
				}

				select {
				case valueStream <- v:
				case <-done:
				}
			}
		}
	}()

	return valueStream
}
