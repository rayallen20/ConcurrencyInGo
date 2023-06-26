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
