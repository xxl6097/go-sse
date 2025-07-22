package main

import (
	"fmt"
	"github.com/xxl6097/go-sse/pkg/sse"
	"net"
	"net/http"
	"strings"
)

func GetLocalMac() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("获取网络接口失败：", err)
		return ""
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.HardwareAddr != nil {
			devMac := strings.ReplaceAll(iface.HardwareAddr.String(), ":", "")
			//fmt.Println(iface.Name, ":", devMac)
			return devMac
		}
	}
	return ""
}

func GetHostIp() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	//var ips []net.IP
	for _, address := range addrList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.IsPrivate() {
			if ipNet.IP.To4() != nil {
				//ip = ipNet.IP.String()
				//break
				ip := ipNet.IP.To4()
				//fmt.Println(ip[0])
				switch ip[0] {
				case 10:
					return ipNet.IP.String()
				case 192:
					return ipNet.IP.String()
				}
			}
		}
	}
	//fmt.Println(ips)
	return ""
}

func main() {
	mac := GetLocalMac()
	ip := GetHostIp()
	fmt.Println(mac)
	fmt.Println(ip)

	url := "http://uuxia.cn:7001/api/sse"
	sse.NewClient(url).
		BasicAuth("admin", "het002402").
		ListenFunc(func(s string) {
			fmt.Printf("SSE: %s\n", s)
		}).Header(func(header *http.Header) {
		header.Add("Sse-Event-IP-Address", ip)
		header.Add("Sse-Event-MAC-Address", mac)
	}).Done()
	select {}
}
