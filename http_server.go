package main

import (
	"fmt"
	"log"
	"net/http"
)

// 简单的路由和中间件实现

type route struct{}

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

func main() {
	var r route
	//http.HandleFunc("/hello", hello) // 这也为路由实现，如果handleFunc和http.Handler一起存在，前者不生效
	log.Fatalln(http.ListenAndServe(":8000", &r))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}
