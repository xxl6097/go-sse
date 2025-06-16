package sse

import (
	"github.com/xxl6097/go-sse/interval"
	"github.com/xxl6097/go-sse/pkg/sse/iface"
)

func NewSseServer(callback interval.OnSseServer) iface.ISseServer {
	return interval.NewServer(callback)
}
