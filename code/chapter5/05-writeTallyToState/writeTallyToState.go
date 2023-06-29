package main

import "fmt"

func main() {
	result := add(1, 2, 3, 4, 5, 6, 7, 8, 9)
	writeTallyToState(result)
}

func add(values ...int) int {
	result := 0
	for _, value := range values {
		result += value
	}

	return result
}

func writeTallyToState(tally int) {
	// 此处使用fmt.Sprintf来模拟一个写入操作
	content := fmt.Sprintf("Tally is %v\n", tally)
	fmt.Printf(content)
}
