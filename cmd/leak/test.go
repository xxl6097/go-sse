package leak

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/xxl6097/go-sse/pkg/sse"
	"github.com/xxl6097/go-sse/pkg/sse/isse"
	"github.com/xxl6097/go-sse/static"
	"log"
	"net/http"
	"sync"
	"time"
)

// Server 结构体管理 HTTP 服务状态
type Server struct {
	server     *http.Server
	shutdownWG sync.WaitGroup
}

// NewServer 创建并配置 HTTP 服务器
func NewServer(addr string) *Server {
	router := mux.NewRouter()
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 注册路由处理器
	router.HandleFunc("/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	router.HandleFunc("/leaky", leakyHandler) // 故意泄漏的端点

	router.Handle("/favicon.ico", http.FileServer(static.FileSystem)).Methods(http.MethodGet, http.MethodOptions)

	serv := sse.New().
		InvalidateFun(func(request *http.Request) (string, error) {
			return time.Now().Format("20060102150405.999999999"), nil
		}).
		Register(func(server isse.ISseServer, client *isse.Client) {
			//server.Stream("内置丰富的开发模板，包括前后端开发所需的所有工具，如pycharm、idea、navicat、vscode以及XTerminal远程桌面管理工具等模板，用户可以轻松部署和管理各种应用程序和工具", time.Millisecond*500)
		}).
		UnRegister(nil).
		Done(context.Background())
	router.HandleFunc("/events", serv.Handler())

	return &Server{server: srv}
}

// 模拟泄漏的处理器：启动永不退出的 Goroutine
func leakyHandler(w http.ResponseWriter, r *http.Request) {
	//ch := make(chan struct{})
	//go func() {
	//	ch <- struct{}{} // 永久阻塞（无接收方）
	//}()
	w.Write([]byte("Leaky endpoint triggered"))
}

// Start 异步启动 HTTP 服务
func (s *Server) Start() {
	s.shutdownWG.Add(1)
	go func() {
		defer s.shutdownWG.Done()
		log.Printf("Server starting on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
}

// Stop 安全关闭服务（含超时控制）
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
	s.shutdownWG.Wait()
	log.Println("Server stopped")
}
