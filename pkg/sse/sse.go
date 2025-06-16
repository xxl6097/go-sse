package sse

import (
	"github.com/xxl6097/go-sse/internal"
	"github.com/xxl6097/go-sse/pkg/sse/iface"
)

func NewSseServer(callback iface.OnSseServer) iface.ISseServer {
	return internal.NewServer(callback)
}
