package main

import (
	"context"
	"fmt"
)

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func main() {
	ProcessRequest("jane", "abc123")
}

func ProcessRequest(userId, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userId)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"Handling response for %v (%v)\n",
		UserId(ctx),
		AuthToken(ctx),
	)
}

func UserId(ctx context.Context) string {
	return ctx.Value(ctxUserID).(string)
}

func AuthToken(ctx context.Context) string {
	return ctx.Value(ctxAuthToken).(string)
}
