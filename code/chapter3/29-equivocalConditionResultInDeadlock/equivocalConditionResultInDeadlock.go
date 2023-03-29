package main

import "fmt"

func main() {
	stringSteam := make(chan string)
	go func() {
		if 0 != 1 {
			return
		}
		stringSteam <- "Fuck"
	}()

	fmt.Println(<-stringSteam)
}
