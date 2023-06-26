package main

import "fmt"

func main() {
	values := []int{1, 2, 3, 4}
	for _, value := range values {
		fmt.Println(multiply(add(multiply(value, 2), 1), 2))
	}
}

func multiply(value, multiplier int) int {
	return value * multiplier
}

func add(value, additive int) int {
	return value + additive
}
