package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	pd "vehicleLicensePlateRecognitionGateway/grpcProto"
)

type server struct{}

//打招呼的服务
//func(对象)函数名（context,客户端发送过来的参数）(返回给客户端的参数,错误)
func (this *server) Sayhello(ctx context.Context, in *pd.HelloReq) (out *pd.HelloRsp, err error) {
	return &pd.HelloRsp{Msg: in.Msg + ";你也好，我是服务端"}, err
}

//说名字的服务
func (this *server) Sayname(ctx context.Context, in *pd.NameReq) (out *pd.NameRsp, err error) {
	return &pd.NameRsp{Name: in.Name + ";我的名字叫服务端10010"}, err
}

func Server() {

	//1、创建网络
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("网络错误", err)
	}

	//2、创建grpc的服务
	ser := grpc.NewServer()

	//3、注册服务
	pd.RegisterHelloServerServer(ser, &server{})

	//4、等待网络连接
	fmt.Println("等待连接中.....")
	err = ser.Serve(listener)
	if err != nil {
		fmt.Println("网络错误", err)
	}

}
