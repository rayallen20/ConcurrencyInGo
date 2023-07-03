package client

import "context"

type APIConnection struct{}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	// 此处假装在调用读取文件的API
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	// 此处假装在调用解析地址为IP的API
	return nil
}

func Open() *APIConnection {
	return &APIConnection{}
}
