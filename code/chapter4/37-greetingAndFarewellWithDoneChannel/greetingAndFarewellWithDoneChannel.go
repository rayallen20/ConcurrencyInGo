package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}()

	wg.Wait()
}

func printGreeting(done <-chan interface{}) error {
	greeting, err := genGreeting(done)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(done <-chan interface{}) error {
	farewell, err := genFarewell(done)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genFarewell(done <-chan interface{}) (string, error) {
	switch language, err := locale(done); {
	case err != nil:
		return "", err
	case language == "EN/US":
		return "Goodbye", nil
	}
	return "", errors.New("unsupported language")
}

func genGreeting(done <-chan interface{}) (string, error) {
	switch language, err := locale(done); {
	case err != nil:
		return "", err
	case language == "EN/US":
		return "Hello", nil
	}
	return "", errors.New("unsupported language")
}

func locale(done <-chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("canceled")
	case <-time.After(1 * time.Second):
	}
	return "EN/US", nil
}
