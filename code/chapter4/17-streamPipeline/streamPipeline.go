package main

import "fmt"

func main() {
	multiply := func(value, multiplier int) int {
		return value * multiplier
	}

	add := func(value, additive int) int {
		return value + additive
	}

	values := []int{1, 2, 3, 4}
	for _, value := range values {
		fmt.Println(multiply(add(multiply(value, 2), 1), 2))
	}
}
