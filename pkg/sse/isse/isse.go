package isse

import (
	"net/http"
	"time"
)

type InvalidateType func(*http.Request) (string, error)
type ClientType func(ISseServer, *Client)

// Event 表示一个 SSE 事件
type Event struct {
	ID      string      `json:"id,omitempty"`
	Event   string      `json:"event,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// Client 表示一个客户端连接
type Client struct {
	ID         string        `json:"id,omitempty"`
	GroupID    string        `json:"groupId,omitempty"`
	IpAddress  string        `json:"ipAddress,omitempty"`
	MacAddress string        `json:"macAddress,omitempty"`
	SendChan   chan Event    `json:"-"`
	CloseChan  chan struct{} `json:"-"`
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
