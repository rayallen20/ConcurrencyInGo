package multiLimiterI

import (
	"context"
	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
	Limit() rate.Limit
}
