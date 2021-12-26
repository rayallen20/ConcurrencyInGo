# sync.Once基础

## PART1. 功能描述

`sync.Once`是GO语言提供的一种使函数只执行1次的对象实现,作用和`init()`函数类似,但有所不同:

1. `init()`函数是在包首次被import时执行,且只执行1次
2. `sync.Once`是在代码运行中需要的时候执行,且只执行1次

由此可知,`sync.Once`的使用场景为:**当一个函数不希望在程序一开始就被执行时,可以使用`sync.Once`**

## PART2. 例子

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once

	onceFunc := func() {
		fmt.Println("Only once")
	}

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceFunc)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
```

运行结果:

```
go run onceDemo1.go 
Only once
```

可以看到,虽然`once.Do(onceFunc)`循环了10次,但匿名函数`onceFunc`却只执行了1次.

## PART3. 源码

`sync.Once`

```go
// Once is an object that will perform exactly one action.
// Once是一个用于精确地只执行1次函数的对象
//
// A Once must not be copied after first use.
type Once struct {
	// done indicates whether the action has been performed.
	// It is first in the struct because it is used in the hot path.
	// The hot path is inlined at every call site.
	// Placing done first allows more compact instructions on some architectures (amd64/386),
	// and fewer instructions (to calculate offset) on other architectures.
	
	// done表示一个函数是否被执行过.把该字段作为结构体的第1个字段是因为该字段经常被热点路径使用.
	done uint32
	m    Mutex
}
```



