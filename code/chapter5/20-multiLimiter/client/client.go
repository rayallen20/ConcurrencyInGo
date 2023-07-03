package client

import (
	"code/chapter5/20-multiLimiter/multiLimiter"
	"code/chapter5/20-multiLimiter/multiLimiterI"
	"context"
	"golang.org/x/time/rate"
	"time"
)

type APIConnection struct {
	rateLimiter multiLimiterI.RateLimiter
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}

	// 此处假装在调用读取文件的API
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}

	// 此处假装在调用解析地址为IP的API
	return nil
}

func Open() *APIConnection {
	// 秒级限速器 限制每秒最多2个事件 桶深度为1
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)

	// 分级限速器 限制每分钟最多10个事件 桶深度为10
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10)

	return &APIConnection{
		rateLimiter: multiLimiter.NewMultiLimiter(secondLimit, minuteLimit),
	}
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}
