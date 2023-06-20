# context

## 1. `WithCancel()`

- `context.Background()`:返回一个根节点
- `context.WithCancel(parent Context) (ctx Context, cancel CancelFunc)`:返回一个parent的副本,该副本具有新的`Done` channel.副本的`Done` channel在调用`cancel()`函数时或parent节点的父节点的`Done` channel关闭时关闭,以先发生的为准.取消该context会释放与其相关的资源,因此在代码中,当此上下文中的操作完成后,应尽快调用`cancel()`.**即:调用`cancel()`时,取消该`cancel()`对应的ctx及其所有子ctx的运行**.

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	wg.Add(1)
	go doSelect(ctx)
	time.Sleep(3 * time.Second)
	cancel()
	wg.Wait()
}

func doSelect(ctx context.Context) {
LOOP:
	for {
		fmt.Printf("select data from DB\n")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LOOP
		default:
			continue
		}
	}
	wg.Done()
}
```

运行结果:

```
go run withCancel.go 
select data from DB
select data from DB
select data from DB
```

## 2. `WithDeadline()`

- `context.WithDeadline(parent Context, d time.Time) (Context, CancelFunc)`:返回一个parent的副本,该副本的deadline不会晚于`d`.若父context的deadline已经早于`d`,则`WithDeadline(parent, d)`在语义上等同父context.返回的context的`Done` channel在到达deadline时、调用返回的`cancel()`函数时或父context的`Done` channel被关闭时关闭,以这三者中先发生的为准.取消该context会释放与其相关的资源,因此在代码中,当此上下文中的操作完成后,应尽快调用`cancel()`.**即:调用`cancel()`时,取消该`cancel()`对应的ctx及其所有子ctx的运行;若不显式调用`cancel()`,则当时间到达`d`时,由`WithDeadline()`调用`cancel()`**.

	- 注:`context.WithDeadline()`相当于接收了一个绝对时刻,当到达这个时刻时即调用`cancel()`

##### 当绝对时刻到达时由`context.WithDeadline()`调用`cancel()`

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	ctx := context.Background()
	ctx, _ = context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	wg.Add(1)
	go doSelect(ctx)
	wg.Wait()
}

func doSelect(ctx context.Context) {
LOOP:
	for {
		fmt.Printf("select data from DB\n")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LOOP
		default:
			continue
		}
	}
	wg.Done()
}
```

运行结果:

```
go run withDeadline.go 
select data from DB
select data from DB
select data from DB
```

##### 绝对时刻尚未到达时显式调用`cancel()`

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	wg.Add(1)
	go doSelect(ctx)
	time.Sleep(2 * time.Second)
	cancel()
	wg.Wait()
}

func doSelect(ctx context.Context) {
LOOP:
	for {
		fmt.Printf("select data from DB\n")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LOOP
		default:
			continue
		}
	}
	wg.Done()
}
```

运行结果:

```
go run withDeadlineCallCancel.go 
select data from DB
select data from DB
```

## 3. `context.WithTimeout()`

- `WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)`:本质上还是`context.WithTimeout()`,只是入参变成了一个相对时段,context中的操作执行时长到达该时段长度时,自动调用`cancel()`

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}
```

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 3*time.Second)
	wg.Add(1)
	go doSelect(ctx)
	wg.Wait()
}

func doSelect(ctx context.Context) {
LOOP:
	for {
		fmt.Printf("select data from DB\n")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LOOP
		default:
			continue
		}
	}
	wg.Done()
}
```

运行结果:

```
go run withTimeout.go 
select data from DB
select data from DB
select data from DB
```

## 4. `context.WithValue()`

- `WithValue(parent Context, key, val any) Context`:返回一个parent的副本,其中与key关联的值为val.仅将上下文范围内的值用于在进程间或API间传递,而不是将一个可选参数传递给函数.填入的key必须是可比较的,不应该是字符串或其他任何内置类型,以避免包之间的冲突.`WithValue()`的使用者应该为key单独定义类型.为了避免给空接口分配内存,key通常是具体类型的结构体实例,或者key的类型是一个指针或接口.注意,上下文的值是传递给后续函数调用的,因此在使用`WithValue()`时要小心不要将敏感或大量数据存储在上下文中,以避免不必要的内存使用.

```go
package main

import (
	"context"
	"fmt"
	"sync"
)

type User struct {
	Id   int
	Name string
}

var wg sync.WaitGroup

func main() {
	user := User{Id: 1, Name: "张三"}
	parent := context.Background()
	ctx := context.WithValue(parent, "user", user)
	ctx = context.WithValue(ctx, "requestId", generateRequestID())
	wg.Add(1)
	go work(ctx)
	wg.Wait()
}

// generateRequestID 生成并返回一个唯一的请求ID
func generateRequestID() string {
	return "12345"
}

func work(ctx context.Context) {
	user := ctx.Value("user").(User)
	requestId := ctx.Value("requestId").(string)
	fmt.Printf("requestId = %s, user = %+v\n", requestId, user)
	wg.Done()
}
```

运行结果:

```
go run withValue.go 
requestId = 12345, user = {Id:1 Name:张三}
```