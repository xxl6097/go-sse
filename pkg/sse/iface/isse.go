package iface

import (
	"net/http"
	"time"
)

type InvalidateType func(*http.Request) (string, error)
type ClientType func(ISseServer, *Client)

// Event 表示一个 SSE 事件
type Event struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Data  string `json:"data"`
}

// Client 表示一个客户端连接
type Client struct {
	ID        string        `json:"id"`
	GroupID   string        `json:"groupId"`
	SendChan  chan Event    `json:"-"`
	CloseChan chan struct{} `json:"-"`
}

//type OnSseServer interface {
//	OnRegister(*Client)
//	OnUnRegister(*Client)
//	Invalidate(*http.Request) (bool, string)
//}

type ISseServer interface {
	GetClients() map[string]*Client
	Stream(response string, interval time.Duration)
	Broadcast(event Event)
	Send(event Event)
	SendToGroup(event Event)
	Handler() http.HandlerFunc
}
