package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	salutation := "Hello"

	wg.Add(1)

	go func() {
		defer wg.Done()
		salutation = "Welcome"
	}()

	wg.Wait()

	fmt.Println(salutation)
}
