package main

import (
	"io"
	"log"
	"net"
)

const dst_host = "www.baidu.com:80"

// 假设 用户和dst不通，但是部署这个服务的设备和dst通，将这台设备收到的内容发送给dst，将dst返回的内容，返回给客户

func main() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal("错误", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("accept出错", err)
		}
		handler(conn)
	}
}

func handler(src net.Conn) {
	dst, err := net.Dial("tcp", dst_host)
	if err != nil {
		log.Fatalln("连接出错", err)
	}
	defer dst.Close()

	// 防止阻塞，使用协程
	go func() {
		// 将源的内容发送给dst
		_, err := io.Copy(dst, src)
		if err != nil {
			log.Fatalln("源到目的copy出错", err)
		}
	}()

	_, err = io.Copy(src, dst)

	if err != nil {
		log.Fatalln("目的到源copy出错", err)
	}
}
