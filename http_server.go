package main

import (
	"fmt"
	"log"
	"net/http"
)

// 简单的路由和中间件实现

type route struct{}

// 构建路由函数
func (r route)ServeHTTP(w http.ResponseWriter,req *http.Request){
	switch req.URL.Path {
	case "/a":
		fmt.Fprint(w,"is a")
	case "/b":
		fmt.Fprintf(w,"is b")
	default:
		http.Error(w,"not found",http.StatusNotFound)
	}

}

func main () {
	var r route
	log.Fatalln(http.ListenAndServe(":8000",&r))
}