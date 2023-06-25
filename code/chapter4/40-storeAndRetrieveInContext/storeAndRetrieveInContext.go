package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest("jane", "abc123")
}

func ProcessRequest(userId, authToken string) {
	ctx := context.WithValue(context.Background(), "userId", userId)
	ctx = context.WithValue(ctx, "authToken", authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	userId := ctx.Value("userId").(string)
	authToken := ctx.Value("authToken").(string)

	fmt.Printf("Handling response for %v (%v)\n", userId, authToken)
}
