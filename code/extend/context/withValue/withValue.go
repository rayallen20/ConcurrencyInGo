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
