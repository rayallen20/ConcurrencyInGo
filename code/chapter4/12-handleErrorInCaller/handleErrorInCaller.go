package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Error    error
	Response *http.Response
}

func main() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{
					Error:    err,
					Response: resp,
				}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()

		return results
	}

	done := make(chan interface{})
	defer close(done)

	errCounter := 0
	urls := []string{"https://www.baidu.com", "http://badHostFoo", "http://badHostBar", "http://badHostBaz"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCounter++
			if errCounter >= 3 {
				fmt.Printf("Occur too many errors, breaking!\n")
				break
			}
			continue
		}

		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
