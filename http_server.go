package main

import (
	"fmt"
	"log"
	"net/http"
)

// 简单的路由和中间件实现

type route struct{}

type logger struct {
	Inner http.Handler
}

// 构建路由函数
func (r route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/a":
		fmt.Fprint(w, "is a")
	case "/b":
		fmt.Fprintf(w, "is b")
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}

}

// 两个serveHttp，先执行这个，这里再调用真正的servehttp
func (l logger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(11)
	l.Inner.ServeHTTP(w, req)
	fmt.Println(22)
}

func main() {
	var r route
	//http.HandleFunc("/hello", hello) // 这也为路由实现，如果handleFunc和http.Handler一起存在，前者不生效
	l := logger{Inner: r}
	log.Fatalln(http.ListenAndServe(":8000", &l))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}
