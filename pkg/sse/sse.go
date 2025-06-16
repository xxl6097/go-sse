package sse

import (
	"github.com/xxl6097/go-sse/internal"
	"github.com/xxl6097/go-sse/pkg/sse/iface"
)

type sseserver struct {
	server *internal.Server
}

func New() *sseserver {
	return &sseserver{
		server: internal.NewServer(),
	}
}

func (s *sseserver) InvalidateFun(fn iface.InvalidateType) *sseserver {
	s.server.InvalidateFun(fn)
	return s
}

func (s *sseserver) Register(fn iface.ClientType) *sseserver {
	s.server.Register(fn)
	return s
}

func (s *sseserver) UnRegister(fn iface.ClientType) *sseserver {
	s.server.UnRegister(fn)
	return s
}

func (s *sseserver) Done() iface.ISseServer {
	return s.server.Done()
}
