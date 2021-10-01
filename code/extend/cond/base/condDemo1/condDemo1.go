package main

import (
	"fmt"
	"sync"
	"time"
)

var mail = make(chan string)

func main() {
	go func() {
		<- mail
		fmt.Println("get chance to do something")
	}()

	time.Sleep(5 * time.Second)
	mail <- "this is a chance to do something"
	time.Sleep(2 * time.Second)

	var l sync.Mutex
	cond := sync.NewCond(&l)
	cond.Wait()
}
