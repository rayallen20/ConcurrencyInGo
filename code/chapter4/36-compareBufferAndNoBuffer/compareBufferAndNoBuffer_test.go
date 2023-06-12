package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// BenchmarkUnbufferedWrite 直接向文件写入
func BenchmarkUnbufferedWrite(b *testing.B) {
	performWrite(b, tmpFileOrFatal())
}

// BenchmarkBufferedWrite 使用带缓冲的writer向文件写入
func BenchmarkBufferedWrite(b *testing.B) {
	// 创建一个带缓冲的writer
	bufferedFile := bufio.NewWriter(tmpFileOrFatal())
	performWrite(b, bufferedFile)
}

// tmpFileOrFatal 创建临时文件
func tmpFileOrFatal() *os.File {
	file, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	return file
}

// performWrite 向给定的writer中写入
func performWrite(b *testing.B, writer io.Writer) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for bt := range take(done, repeat(done, byte(0)), b.N) {
		writer.Write([]byte{bt.(byte)})
	}
}

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)
		for {
			for _, value := range values {
				select {
				case <-done:
					return
				case valueStream <- value:
				}
			}
		}
	}()

	return valueStream
}

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}
