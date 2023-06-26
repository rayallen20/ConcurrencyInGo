package main

func main() {
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
