package main

import (
	"flag"
	"fmt"
	"goshop_srvs/user_srv/handler"
	"goshop_srvs/user_srv/proto"
	"net"

	"google.golang.org/grpc"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.String("port", "50051", "端口号")

	flag.Parse()
	fmt.Println("ip:", *IP)
	fmt.Println("port:", *Port)

	server := grpc.NewServer()
	// 注册
	proto.RegisterUserServer(server, &handler.UserServer{})
	// 启动
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	err = server.Serve(lis) // 没有直接结束连接
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}

/*
 cd到 这个文件
 执行 go build mian.go
 会生成 main.exe

 main.exe -h 会出现命令使用帮助
    -ip String  ip地址
    -port int   端口号

 main.exe -port 50053  启动的时候端口号就变了
*/
