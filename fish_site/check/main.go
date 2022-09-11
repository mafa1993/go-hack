package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// 利用中间件进行身份验证
// 用的第三方包
//go get github.com/gorilla/mux
//go get github.com/urfave/negroni

type badAuth struct {
	Username string
	Password string
}

func (b *badAuth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username != b.Username || password != b.Password {
		http.Error(w, "验证失败", http.StatusAccepted)
		return
	}

	// 在上下文中设置一个新变量username 值为变量username
	ctx := context.WithValue(r.Context(), "username", username)
	r = r.WithContext(ctx)
	next(w, r)
}

func hello(w http.ResponseWriter,r *http.Request){
	username := r.Context().Value("username").(string)  // 获取上下文中设置的username

	fmt.Fprintln(w,"HI"+username)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/hello",hello).Method("GET") // 设置路由，hello

	n := negroni.Classic()
	n.Use(&badAuth{  // 将 hegroni.Handler 实例传递进来
		Username: "admin",
		Password: "password",
	})

	n.UseHandler(r)
	http.ListenAndServe(":8000",n)
}

