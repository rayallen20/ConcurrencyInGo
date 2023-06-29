package main

import "time"

func main() {

}

func runInGoroutine(done <-chan interface{}, valueStream <-chan interface{}) {
	var value interface{}
	resultStream := make(chan interface{})

	select {
	case <-done:
		return
	case value = <-valueStream:
	}

	result := reallyLongCalculation(done, value)

	select {
	case <-done:
		return
	case resultStream <- result:
	}
}

func reallyLongCalculation(done <-chan interface{}, value interface{}) interface{} {
	intermediateResult := longCalculation(done, value)
	return longCalculation(done, intermediateResult)
}

func longCalculation(done <-chan interface{}, value interface{}) interface{} {
	select {
	case <-done:
		return nil
	default:
	}
	time.Sleep(3 * time.Second)
	return value
}
