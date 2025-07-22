package sse

import (
	"bufio"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	URL                string
	username, password string
	fn                 func(string)
	headerFn           func(*http.Header)
	done               chan struct{}
}

func NewClient(url string) *Client {
	return &Client{URL: url, done: make(chan struct{})}
}

func (c *Client) BasicAuth(username, password string) *Client {
	c.username = username
	c.password = password
	return c
}

func (c *Client) ListenFunc(fn func(string)) *Client {
	c.fn = fn
	return c
}

func (c *Client) Header(fn func(*http.Header)) *Client {
	c.headerFn = fn
	return c
}

func (c *Client) Done() *Client {
	go c.connect()
	return c
}

func (c *Client) Close() {
	c.done <- struct{}{}
}

func (c *Client) connect() {
	for {
		select {
		case <-c.done:
			return
		default:
			req, _ := http.NewRequest("GET", c.URL, nil)
			req.SetBasicAuth(c.username, c.password)
			if c.headerFn != nil {
				c.headerFn(&req.Header)
			}
			req.Header.Set("Accept", "text/event-stream")
			client := &http.Client{Timeout: 0} // 无超时限制
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("连接失败: %v，5秒后重试", err)
				time.Sleep(5 * time.Second)
				continue
			}
			scanner := bufio.NewScanner(resp.Body)
			var eventType, data string
			for scanner.Scan() {
				line := scanner.Text()
				switch {
				case strings.HasPrefix(line, "data:"):
					data = strings.TrimSpace(line[5:])
					if c.fn != nil {
						c.fn(data)
					}
				case strings.HasPrefix(line, "event:"):
					eventType = strings.TrimSpace(line[6:])
				case line == "" && data != "":
					log.Printf("[%s] %s", eventType, data)
					data, eventType = "", ""
				}
			}
			resp.Body.Close()
		}
	}
}
