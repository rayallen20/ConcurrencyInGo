package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	var stdoutBuffer bytes.Buffer
	defer stdoutBuffer.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintf(&stdoutBuffer, "Producer done.\n")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuffer, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Fprintf(&stdoutBuffer, "Received: %d\n", integer)
	}
}
