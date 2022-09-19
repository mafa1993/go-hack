# 准备工作

1. 使用grpc和protobuf构建
2. 包含三个部分，客户端、服务端、管理程序
3. 创建一个植入程序，定期轮询客户机，并将输出返给服务端
4. go get -u google.golang.org/grpc
5. protoc 命令 protoc -I . --go-grpc_out=.  --go_out=. implant.proto  生成proto接口