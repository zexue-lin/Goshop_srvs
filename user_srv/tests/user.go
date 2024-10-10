package main

import (
	"context"
	"fmt"
	"goshop_srvs/user_srv/proto"

	"google.golang.org/grpc"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

// client是通用的
func Init() {
	// 拨号连接
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	// 远程调用
	userClient = proto.NewUserClient(conn)
	// r, err := c.SayHello(context.Background(), &protoc.HelloRequest{Name: "Jerry"}) // 指针类型
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(r.Message)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.Password)
		checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}


}

func main() {
	Init()
	TestGetUserList()
	conn.Close()
}
