package process

import (
	"code/chapter4/42-storeAndRetrieveInContextWithCustomKey/ctxCustomKey/response"
	"context"
)

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func ProcessRequest(userId, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userId)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	response.HandleResponse(ctx)
}

func UserId(ctx context.Context) string {
	return ctx.Value(ctxUserID).(string)
}

func AuthToken(ctx context.Context) string {
	return ctx.Value(ctxAuthToken).(string)
}
