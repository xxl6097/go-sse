package sse

import (
	"context"
	"github.com/xxl6097/go-sse/internal"
	"github.com/xxl6097/go-sse/pkg/sse/isse"
)

type sseserver struct {
	server *internal.Server
}

func New() *sseserver {
	return &sseserver{
		server: internal.NewServer(),
	}
}

func (s *sseserver) InvalidateFun(fn isse.InvalidateType) *sseserver {
	s.server.InvalidateFun(fn)
	return s
}

func (s *sseserver) Register(fn isse.ClientType) *sseserver {
	s.server.Register(fn)
	return s
}

func (s *sseserver) UnRegister(fn isse.ClientType) *sseserver {
	s.server.UnRegister(fn)
	return s
}

func (s *sseserver) Done(ctx context.Context) isse.ISseServer {
	return s.server.Done(ctx)
}
