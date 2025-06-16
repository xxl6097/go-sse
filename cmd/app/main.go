package main

import (
	"errors"
	"github.com/xxl6097/go-sse/pkg/sse"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//// 创建 SSE 服务器
	//sseServer := sse.NewServer()
	//sseServer.Start()
	//
	//// 设置路由
	//http.HandleFunc("/events", sseServer.SubscribeHandler())
	//
	//// 模拟定时发送消息
	//go func() {
	//	counter := 0
	//	ticker := time.NewTicker(5 * time.Second)
	//	defer ticker.Stop()
	//
	//	for range ticker.C {
	//		counter++
	//		sseServer.Broadcast(sse.Event{
	//			ID:    time.Now().Format(time.RFC3339Nano),
	//			Event: "message",
	//			Data:  "This is message #" + string(counter),
	//		})
	//	}
	//}()

	//broker := interval.NewAdvancedBroker(5)
	//broker.StreamResponse("你好，你是谁，从哪里来，年龄，性别，姓名，一一报来！", time.Second)
	//http.Handle("/events", broker)

	serv := sse.New().
		InvalidateFun(func(request *http.Request) (bool, string) {
			return true, time.Now().Format("20060102150405.999999999")
		}).
		Register(nil).
		UnRegister(nil).
		Done()
	http.HandleFunc("/events", serv.Handler())

	// 提供静态文件服务
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// 启动服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	log.Println("Server started on :8080")

	// 优雅关闭
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 优雅关闭服务器
	if err := server.Shutdown(nil); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
