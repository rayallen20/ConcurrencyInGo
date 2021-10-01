package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

func main() {
	producer := func(wg *sync.WaitGroup, lock sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			lock.Lock()
			lock.Unlock()
			time.Sleep(1)
		}
	}

	observer := func(wg *sync.WaitGroup, lock sync.Locker) {
		defer wg.Done()
		lock.Lock()
		defer lock.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)

		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}
		wg.Wait()

		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	defer tw.Flush()

	var rwMutex sync.RWMutex
	fmt.Fprintf(tw, "Readers\tMutex\tRWMutex\n")

	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &rwMutex, rwMutex.RLocker()),
			test(count, &rwMutex, &rwMutex),
			)
	}
}
