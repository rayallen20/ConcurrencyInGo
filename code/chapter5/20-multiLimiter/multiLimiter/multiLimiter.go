package multiLimiter

import (
	"code/chapter5/20-multiLimiter/multiLimiterI"
	"context"
	"golang.org/x/time/rate"
	"sort"
)

type MultiLimiter struct {
	limiters []multiLimiterI.RateLimiter
}

func (m *MultiLimiter) Wait(ctx context.Context) error {
	for _, limiter := range m.limiters {
		if err := limiter.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (m *MultiLimiter) Limit() rate.Limit {
	return m.limiters[0].Limit()
}

func NewMultiLimiter(limiters ...multiLimiterI.RateLimiter) *MultiLimiter {
	byLimit := func(i, j int) bool {
		// 按限流器限制的速度从低到高排序
		return limiters[i].Limit() < limiters[j].Limit()
	}

	sort.Slice(limiters, byLimit)
	return &MultiLimiter{limiters: limiters}
}
