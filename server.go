package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"raspberrypi/protobuf"
)

type server struct {
	protobuf.UnimplementedInteractServer
}

func (s *server) SetAddr(ctx context.Context, setAddrRequest *protobuf.SetAddrRequest) (*protobuf.SetAddrResponse, error) {
	fmt.Printf("%s:%s\n", setAddrRequest.InterfaceName, setAddrRequest.InterfaceAddr)
	return &protobuf.SetAddrResponse{
		Reply: "success",
	}, nil
}

func main() {
	err := serverCore()
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func serverCore() error {
	listeningPort := 8972
	// 监听地址
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", listeningPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()                         // 创建gRPC服务器
	protobuf.RegisterInteractServer(grpcServer, &server{}) // 在gRPC服务端注册服务
	err = grpcServer.Serve(lis)                            // 启动服务
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
