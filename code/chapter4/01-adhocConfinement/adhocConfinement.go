package main

import "fmt"

func main() {
	data := make([]int, 4)

	handleData := make(chan int)
	go loopData(handleData, data)

	for num := range handleData {
		fmt.Printf("%d\n", num)
	}
}

func loopData(handleData chan<- int, data []int) {
	defer close(handleData)
	for i := range data {
		handleData <- data[i]
	}
}
