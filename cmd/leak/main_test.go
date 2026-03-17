package leak

import (
	"go.uber.org/goleak"
	"net/http"
	"testing"
)

// go test -v -run TestLeakyEndpoint

func TestHealthyEndpoint(t *testing.T) {
	defer goleak.VerifyNone(t) // 泄漏检测

	// 启动测试服务器
	server := NewServer(":8001")
	server.Start()
	defer server.Stop()

	// 发送健康检查请求
	//resp, err := http.Get("http://localhost" + server.server.Addr + "/healthy")
	resp, err := http.Get("http://localhost:8081/mqtt/test")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m) // 全局泄漏检测
}
