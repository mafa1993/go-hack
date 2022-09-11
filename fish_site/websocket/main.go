package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

//"github.com/gorilla/websocket"

var (
	upgrader = websocket.Upgrader{
		// 设置websocket的源白名单
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	listenAddr string // http服务地址
	wsAddr     string // websocket 服务地址
	jsTemplate *template.Template
)

func init() {
	flag.StringVar(&listenAddr, "listen-addr", "", "http服务地址")
	flag.StringVar(&wsAddr, "ws-addr", "", "ws服务地址")
	flag.Parse()
	log.Println(wsAddr)
	log.Println(listenAddr)
	var err error
	jsTemplate, err = template.ParseFiles("logger.js")
	if err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", serveWS)
	r.HandleFunc("/k.js", serveFile)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./"))) // 静态文件服务器创建
	log.Fatal(http.ListenAndServe(":8080", r))
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "建立ws错误", http.StatusBadRequest)
	}

	defer conn.Close()

	fmt.Printf("conn from %s\n", conn.RemoteAddr().String())

	for {
		_, msg, err := conn.ReadMessage() // 读取ws消息
		if err != nil {
			return
		}
		log.Println(msg)

		fmt.Printf("From %s,%s\n", conn.RemoteAddr().String(), string(msg))
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript") // 设置类型为js
	jsTemplate.Execute(w, wsAddr)
}
