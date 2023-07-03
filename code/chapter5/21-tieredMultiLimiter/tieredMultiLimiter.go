package main

import (
	"code/chapter5/21-tieredMultiLimiter/client"
	"context"
	"log"
	"os"
	"sync"
)

func main() {
	defer log.Printf("Done.\n")

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := client.Open()

	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot read file: %v\n", err)
			}

			log.Printf("ReadFile\n")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot resolve address: %v\n", err)
			}

			log.Printf("ResolveAddress\n")
		}()
	}

	wg.Wait()
}
