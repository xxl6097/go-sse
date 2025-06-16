package iface

import "net/http"

// Event 表示一个 SSE 事件
type Event struct {
	ID    string
	Event string
	Data  string
}

type ISseServer interface {
	Broadcast(event Event)
	Send(event Event)
	SendToGroup(event Event)
	Handler() http.HandlerFunc
}
