package main

import (
	"fmt"
	"net/http"
)

func main() {
	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.baidu.com", "https://badHost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}
}

func checkStatus(done <-chan interface{}, urls ...string) <-chan *http.Response {
	responses := make(chan *http.Response)
	go func() {
		defer close(responses)
		for _, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				continue
			}

			select {
			case <-done:
				return
			case responses <- resp:
			}
		}
	}()

	return responses
}
