package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pd "vehicleLicensePlateRecognitionGateway/grpcProto"
)

func Client() {

	//客户端连接服务器
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("网络异常", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	//获取grpc句柄
	c := pd.NewHelloServerClient(conn)

	//通过句柄调用函数
	req1, err := c.Sayhello(context.Background(), &pd.HelloReq{Msg: "你好，我是客户端"})
	if err != nil {
		fmt.Println("Sayhello服务调用失败", err)
	}
	fmt.Println("调用sayhello函数:", req1.Msg)

	req2, err := c.Sayname(context.Background(), &pd.NameReq{Name: "我名字叫客户端10086"})
	if err != nil {
		fmt.Println("Sayname服务调用失败", err)
	}
	fmt.Println("调用sayname函数:", req2.Name)
}
