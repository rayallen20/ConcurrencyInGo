package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 1*time.Second)
	deadline, _ := ctx.Deadline()
	fmt.Printf("deadline = %#v\n", deadline.String())

	target := time.Now().Add(1 * time.Minute)
	fmt.Printf("target = %#v\n", target.String())

	sub := deadline.Sub(target).Seconds()
	fmt.Printf("sub = %#v\n", sub)
}
