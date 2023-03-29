package main

import "fmt"

func main() {
	var dataStream chan<- interface{}
	dataStream = make(chan<- interface{})
	fmt.Printf("%#v\n", dataStream)
}
