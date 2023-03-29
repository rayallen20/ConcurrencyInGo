package main

import "fmt"

func main() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Fuck World"
	}()
	fmt.Printf("receive from channel: %s\n", <-stringStream)
}
