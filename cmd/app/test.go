package main

import (
	"github.com/xxl6097/go-sse/pkg/sse/iface"
	"log"
	"net/http"
	"time"
)

// generateClientID 生成唯一的客户端 ID
func generateClientID() string {
	return time.Now().Format("20060102150405.999999999")
}

type SSEServer struct{}

func (S SSEServer) OnRegister(client *iface.Client) {
	//TODO implement me
	log.Printf("iface.Client %s connected", client.ID)
}

func (S SSEServer) OnUnRegister(client *iface.Client) {
	//TODO implement me
	log.Printf("iface.Client %s disconnected", client.ID)
}

func (S SSEServer) Invalidate(request *http.Request) (bool, string) {
	//TODO implement me
	return true, generateClientID()
}

func NewSSEServer() iface.OnSseServer {
	return SSEServer{}
}
