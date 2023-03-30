package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()
		data[0] = 'f'
		var buffer bytes.Buffer
		for _, byteData := range data {
			fmt.Fprintf(&buffer, "%c", byteData)
		}
		fmt.Printf("%s\n", buffer.String())
	}

	var wg sync.WaitGroup
	data := []byte("golang")
	wg.Add(2)
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])
	wg.Wait()

	fmt.Printf("%v\n", data)
}
