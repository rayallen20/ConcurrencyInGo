package storeAndRetrive

import (
	"context"
	"fmt"
)

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"Handling response for %v (%v)\n",
		UserId(ctx),
		AuthToken(ctx),
	)
}
