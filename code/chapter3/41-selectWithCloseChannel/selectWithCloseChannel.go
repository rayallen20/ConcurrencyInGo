package main

import "fmt"

func main() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Counter, c2Counter int

	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Counter++
		case <-c2:
			c2Counter++
		}
	}

	fmt.Printf("c1Counter = %d, c2Counter = %d\n", c1Counter, c2Counter)
}
