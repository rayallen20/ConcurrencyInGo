package response

import (
	"code/chapter4/42-storeAndRetrieveInContextWithCustomKey/ctxCustomKey/process"
	"context"
	"fmt"
)

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"Handling response for %v (%v)\n",
		process.UserId(ctx),
		process.AuthToken(ctx),
	)
}
