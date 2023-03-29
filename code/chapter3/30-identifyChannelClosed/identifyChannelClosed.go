package main

import "fmt"

func main() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Fuck"
	}()
	salutation, ok := <-stringStream
	fmt.Printf("(%v): %v\n", ok, salutation)
}
