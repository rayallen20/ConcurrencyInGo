package main

import (
	"golang.org/x/time/rate"
	"time"
)

func main() {

}

// Per 返回一个 rate.Limit 实例.该实例代表在给定的时段duration内有给定数量eventCount个事件发生的前提下
// 每秒允许的最大事件数量
func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}
