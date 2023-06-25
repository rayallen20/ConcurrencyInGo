package storeAndRetrive

import (
	"context"
)

func ProcessRequest(userId, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userId)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func UserId(ctx context.Context) string {
	return ctx.Value(ctxUserID).(string)
}

func AuthToken(ctx context.Context) string {
	return ctx.Value(ctxAuthToken).(string)
}
