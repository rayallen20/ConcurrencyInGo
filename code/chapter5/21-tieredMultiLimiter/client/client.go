package client

import (
	"code/chapter5/21-tieredMultiLimiter/multiLimiter"
	"code/chapter5/21-tieredMultiLimiter/multiLimiterI"
	"context"
	"golang.org/x/time/rate"
	"time"
)

type APIConnection struct {
	networkLimit multiLimiterI.RateLimiter
	diskLimit    multiLimiterI.RateLimiter
	apiLimit     multiLimiterI.RateLimiter
}

// ReadFile 模拟读取文件 前10个请求每秒1个 从第11个开始每6秒1个 第11个请求会在第10个请求等待2s后执行
func (a *APIConnection) ReadFile(ctx context.Context) error {
	err := multiLimiter.NewMultiLimiter(a.apiLimit, a.diskLimit).Wait(ctx)
	if err != nil {
		return err
	}

	// 此处假装在调用读取文件的API
	return nil
}

// ResolveAddress 模拟解析地址为IP 前10个请求每秒2个 从第11个开始每3秒1个 第11个请求会在第10个请求等待1s后执行
func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	err := multiLimiter.NewMultiLimiter(a.apiLimit, a.networkLimit).Wait(ctx)
	if err != nil {
		return err
	}

	// 此处假装在调用解析地址为IP的API
	return nil
}

func Open() *APIConnection {
	return &APIConnection{
		apiLimit: multiLimiter.NewMultiLimiter(
			// API限速器1: 限制每秒最多2次访问 桶深为2
			rate.NewLimiter(Per(2, time.Second), 2),
			// API限速器2: 限制每分钟最多10次访问 桶深为10
			rate.NewLimiter(Per(10, time.Minute), 10),
		),
		diskLimit: multiLimiter.NewMultiLimiter(
			// 磁盘访问限速器1: 限制每秒最多1次访问 桶深为1
			rate.NewLimiter(rate.Limit(1), 1),
		),
		networkLimit: multiLimiter.NewMultiLimiter(
			// 网络访问限速器1: 限制每秒最多3次访问 桶深为3
			rate.NewLimiter(Per(3, time.Second), 3),
		),
	}
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}
