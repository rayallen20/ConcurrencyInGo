package main

import (
	"fmt"
	"sync"
)

func main() {
	var numCalcCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// 用4KB初始化pool
	calcPool.Put(calcPool.New)
	calcPool.Put(calcPool.New)
	calcPool.Put(calcPool.New)
	calcPool.Put(calcPool.New)

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			mem := calcPool.Get()
			// 此处假定在内存中做了某些速度较快的操作
			defer calcPool.Put(mem)
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created\n", numCalcCreated)
}
