package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

func main() {
	var stdoutBuffer bytes.Buffer
	var lock sync.Mutex
	defer stdoutBuffer.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer func() {
			lock.Lock()
			fmt.Fprintf(&stdoutBuffer, "Producer done.\n")
			lock.Unlock()
			close(intStream)
		}()
		for i := 0; i < 5; i++ {
			lock.Lock()
			fmt.Fprintf(&stdoutBuffer, "Sending: %d\n", i)
			lock.Unlock()
			intStream <- i
		}
	}()

	for integer := range intStream {
		lock.Lock()
		fmt.Fprintf(&stdoutBuffer, "Received: %d\n", integer)
		lock.Unlock()
	}
}
