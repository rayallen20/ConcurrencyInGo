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

// warmServiceConnCache 本函数用于在池中预创建10个连接
func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool{
		New: connectToService,
	}

	// 预创建10个连接放入池中
	for i := 0; i < 10; i++ {
		p.Put(p.New)
	}
	return p
}

// startNetworkDaemon 本函数是一个网络处理程序 每次调用本函数 仅允许1个连接
// 但这个连接是从池中取出的 并非是请求时创建的
func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()
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

			serviceConn := connPool.Get()
			fmt.Fprintln(conn, "")
			connPool.Put(serviceConn)
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
