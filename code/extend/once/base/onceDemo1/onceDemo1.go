package main

import (
	"fmt"
	"time"
)

func main() {
	//var once sync.Once
	//
	//onceFunc := func() {
	//	fmt.Println("Only once")
	//}
	//
	//done := make(chan bool)
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		once.Do(onceFunc)
	//		done <- true
	//	}()
	//}
	//
	//for i := 0; i < 10; i++ {
	//	<-done
	//}

	now := time.Now().Format("2006-01-02")
	fmt.Println(now)

	//threeDaysMiddleUnix := int64(3 * 24 * 60 * 60)
	//threeDaysAfterUnix := time.Now().Unix() + threeDaysMiddleUnix
	//threeDaysUnix := time.Unix(threeDaysAfterUnix, 0)
	threeDaysTime := time.Now().Add(-3 * 24 * time.Hour)
	fmt.Println(threeDaysTime.Format("2006-01-02"))
}
