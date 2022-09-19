package main

import (
	"c2/implant/grpcapi"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
)

// 管理客户端，用于向客户端下发指令

func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.AdminClient
	)

	opts = append(opts, grpc.WithInsecure())

	if conn, err = grpc.Dial("127.0.0.1:8081", opts...); err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client = grpcapi.NewAdminClient(conn)
	var cmd = new(grpcapi.Command)
	cmd.In = os.Args[1]
	ctx := context.Background()

	cmd, _ = client.RunCommand(ctx, cmd) // 发送给管理服务器

	fmt.Println(cmd.Out)

}
