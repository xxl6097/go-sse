package internal

import (
	"encoding/json"
	"github.com/xxl6097/go-sse/pkg/sse/isse"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Server 管理所有 SSE 连接和事件广播
type Server struct {
	clients         map[string]*isse.Client
	clientsMutex    sync.RWMutex
	subscribeChan   chan *isse.Client
	invalidateFn    isse.InvalidateType
	registerFn      isse.ClientType
	unregisterFn    isse.ClientType
	unsubscribeChan chan string
	broadcastChan   chan isse.Event
	p2pChan         chan isse.Event
	groupChan       chan isse.Event
}

func (s *Server) Handler() http.HandlerFunc {
	return s.SubscribeHandler()
}

// NewServer 创建一个新的 SSE 服务器实例
func NewServer() *Server {
	this := Server{
		clients:         make(map[string]*isse.Client),
		subscribeChan:   make(chan *isse.Client),
		unsubscribeChan: make(chan string),
		broadcastChan:   make(chan isse.Event),
		p2pChan:         make(chan isse.Event),
		groupChan:       make(chan isse.Event),
		registerFn:      nil,
		invalidateFn:    nil,
		unregisterFn:    nil,
	}
	return &this
}

func (s *Server) Done() isse.ISseServer {
	go s.eventLoop()
	return s
}

func (s *Server) GetClients() map[string]*isse.Client {
	return s.clients
}

// SubscribeHandler 返回一个处理客户端订阅的 HTTP 处理器
func (s *Server) SubscribeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置必要的 HTTP 头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		id := generateClientID()
		// 校验连接合法性
		if s.invalidateFn != nil {
			tempID, err := s.invalidateFn(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			} else {
				id = tempID
			}
		}

		// 创建新客户端
		client := &isse.Client{
			ID:         id,
			GroupID:    r.Header.Get("Sse-Event-GroupID"),
			IpAddress:  r.Header.Get("Sse-Event-IP-Address"),
			MacAddress: r.Header.Get("Sse-Event-MAC-Address"),
			SendChan:   make(chan isse.Event, 100),
			CloseChan:  make(chan struct{}),
		}

		// 注册客户端
		s.subscribeChan <- client

		// 确保客户端断开时注销
		defer func() {
			s.unsubscribeChan <- client.ID
			close(client.CloseChan)
		}()

		// 保持连接打开，发送事件
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		s.sendHeartbeat(w, flusher)

		// 发送心跳以保持连接活跃
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case event, ok := <-client.SendChan:
				if !ok {
					return
				}
				// 发送事件到客户端
				s.sendEvent(w, flusher, event)

			case <-ticker.C:
				// 发送心跳
				s.sendHeartbeat(w, flusher)

			case <-r.Context().Done():
				log.Printf("sse %s Done", client.ID)
				return
			}
		}
	}
}

func (s *Server) InvalidateFun(fn isse.InvalidateType) *Server {
	s.invalidateFn = fn
	return s
}

func (s *Server) Register(fn isse.ClientType) *Server {
	s.registerFn = fn
	return s
}

func (s *Server) UnRegister(fn isse.ClientType) *Server {
	s.unregisterFn = fn
	return s
}

// Broadcast 向所有客户端广播事件
func (s *Server) Broadcast(event isse.Event) {
	s.broadcastChan <- event
}

// Send 向客户端发送数据
func (s *Server) Send(event isse.Event) {
	s.p2pChan <- event
}

// SendToGroup 向一组客户端发送数据
func (s *Server) SendToGroup(event isse.Event) {
	s.groupChan <- event
}

// eventLoop 处理客户端订阅、取消订阅和事件广播
func (s *Server) eventLoop() {
	for {
		select {
		case client := <-s.subscribeChan:
			s.clientsMutex.Lock()
			s.clients[client.ID] = client
			if s.registerFn != nil {
				s.registerFn(s, client)
			}
			s.clientsMutex.Unlock()
			log.Printf("iface.Client %s connected", client.ID)

		case clientID := <-s.unsubscribeChan:
			s.clientsMutex.Lock()
			if client, exists := s.clients[clientID]; exists {
				if s.unregisterFn != nil {
					s.unregisterFn(s, client)
				}
				close(client.SendChan)
				delete(s.clients, clientID)
			}
			s.clientsMutex.Unlock()
			log.Printf("iface.Client %s disconnected", clientID)

		case event := <-s.broadcastChan:
			s.clientsMutex.RLock()
			for _, client := range s.clients {
				select {
				case client.SendChan <- event:
				default:
					// 客户端缓冲区已满，考虑关闭连接或实现退避策略
					log.Printf("iface.Client %s buffer full, dropping event", client.ID)
				}
			}
			s.clientsMutex.RUnlock()

		case event := <-s.p2pChan:
			s.clientsMutex.RLock()
			client := s.clients[event.ID]
			select {
			case client.SendChan <- event:
			default:
				// 客户端缓冲区已满，考虑关闭连接或实现退避策略
				log.Printf("iface.Client %s buffer full, dropping event", client.ID)
			}
			s.clientsMutex.RUnlock()

		case event := <-s.groupChan:
			s.clientsMutex.RLock()
			for _, client := range s.clients {
				if strings.Compare(strings.ToLower(event.ID), strings.ToLower(client.GroupID)) != 0 {
					continue
				}
				select {
				case client.SendChan <- event:
				default:
					// 客户端缓冲区已满，考虑关闭连接或实现退避策略
					log.Printf("iface.Client %s buffer full, dropping event", client.ID)
				}
			}
			s.clientsMutex.RUnlock()

		}
	}
}

// sendEvent 向客户端发送 SSE 事件
func (s *Server) sendEvent(w http.ResponseWriter, flusher http.Flusher, event isse.Event) {
	data, e := json.Marshal(event)
	if e != nil {
		log.Printf("Error marshaling event: %v", e)
		return
	}

	if _, err := w.Write([]byte("data: " + string(data) + "\n\n")); err != nil {
		log.Printf("Write error: %v", err)
		return
	}
	flusher.Flush()
}

// sendHeartbeat 发送心跳消息以保持连接活跃
func (s *Server) sendHeartbeat(w http.ResponseWriter, flusher http.Flusher) {
	if n, err := w.Write([]byte(": heartbeat\n\n")); err != nil {
		log.Printf("Write error: %v", err)
		return
	} else {
		log.Printf("heartbeat: %v", n)
	}
	flusher.Flush()
}

func (b *Server) Stream(response string, interval time.Duration) {
	go func() {
		for i, char := range response {
			event := isse.Event{
				Payload: string(char),
				ID:      strconv.Itoa(i + 1),
				Event:   "message",
			}
			log.Printf("Stream: %v", string(char))
			b.Broadcast(event)
			time.Sleep(interval)
		}
		b.Broadcast(isse.Event{Event: "end"})
	}()
}

// generateClientID 生成唯一的客户端 ID
func generateClientID() string {
	return time.Now().Format("20060102150405.999999999")
}
