package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func main() {

}

// connectToService 本函数用于模拟创建一个到服务连接
// 此处故意让创建连接这个过程消耗较长的时间
func connectToService() interface{} {
	time.Sleep(time.Second)
	return struct{}{}
}

// startNetworkDaemon 本函数是一个网络处理程序 每次调用本函数 仅允许1个连接
func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	// 为简化基准测试 此处每次仅允许1个连接
	wg.Add(1)

	go func() {
		server, err := net.Listen("tcp", "localhost:8090")
		if err != nil {
			log.Fatalf("can not listen: %v\n", err)
		}
		defer server.Close()

		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("can not accept connection: %v\n", err)
				continue
			}

			connectToService()
			fmt.Fprintln(conn, "")
			conn.Close()
		}
	}()

	return &wg
}

func init() {
	daemonStarted := startNetworkDaemon()
	daemonStarted.Wait()
}

func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8090")
		if err != nil {
			b.Fatalf("can not dial host: %v\n", err)
		}

		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("can not read: %v\n", err)
		}
		conn.Close()
	}
}