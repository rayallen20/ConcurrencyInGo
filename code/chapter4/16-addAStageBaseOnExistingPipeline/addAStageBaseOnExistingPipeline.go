package main

import "fmt"

func main() {
	intCollection := []int{1, 2, 3, 4}
	for _, v := range multiply(add(multiply(intCollection, 2), 1), 2) {
		fmt.Println(v)
	}
}

func multiply(values []int, multiplier int) []int {
	multipliedValues := make([]int, len(values))
	for i, value := range values {
		multipliedValues[i] = multiplier * value
	}

	return multipliedValues
}

func add(values []int, additive int) []int {
	addedValues := make([]int, len(values))
	for i, value := range values {
		addedValues[i] = value + additive
	}

	return addedValues
}
