package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"path"
	"raspberrypi/protobuf"
	"raspberrypi/utils/file"
	"raspberrypi/utils/network"
)

type server struct {
	protobuf.UnimplementedInteractServer
}

func (s *server) SetAddr(ctx context.Context, setAddrRequest *protobuf.SetAddrRequest) (*protobuf.NormalResponse, error) {
	fmt.Printf("%s:%s\n", setAddrRequest.InterfaceName, setAddrRequest.InterfaceAddr)
	// ---------------- 核心逻辑 ----------------
	err := network.SetAddr(setAddrRequest.InterfaceName, setAddrRequest.InterfaceAddr)
	if err != nil {
		fmt.Printf("set addr failed: %v\n", err)
		return &protobuf.NormalResponse{
			Reply: "failed",
		}, nil
	}
	// ---------------- 核心逻辑 ----------------
	return &protobuf.NormalResponse{
		Reply: "success",
	}, nil
}

func (s *server) AddRoute(ctx context.Context, addRouteRequest *protobuf.AddRouteRequest) (*protobuf.NormalResponse, error) {
	fmt.Printf("%s:%s\n", addRouteRequest.DestinationNetworkSegment, addRouteRequest.Gateway)
	// ---------------- 核心逻辑 ----------------
	err := network.AddRoute(addRouteRequest.DestinationNetworkSegment, addRouteRequest.Gateway)
	if err != nil {
		fmt.Printf("add route failed: %v\n", err)
		return &protobuf.NormalResponse{
			Reply: "failed",
		}, nil
	}
	// ---------------- 核心逻辑 ----------------
	return &protobuf.NormalResponse{
		Reply: "success",
	}, nil
}

func (s *server) TransmitFile(ctx context.Context, transmitFileRequest *protobuf.TransmitFileRequest) (*protobuf.NormalResponse, error) {
	fmt.Println("transmit file")
	// ---------------- 核心逻辑 ----------------
	directory := path.Dir(transmitFileRequest.DestinationPath)
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Printf("create directory failed: %v\n", err)
		return &protobuf.NormalResponse{
			Reply: "failed",
		}, nil
	}
	err = file.WriteStringIntoFile(transmitFileRequest.DestinationPath, transmitFileRequest.Content)
	if err != nil {
		fmt.Printf("write file failed: %v\n", err)
		return &protobuf.NormalResponse{
			Reply: "failed",
		}, nil
	}
	fmt.Printf("write file %s success\n", transmitFileRequest.DestinationPath)
	// ---------------- 核心逻辑 ----------------
	return &protobuf.NormalResponse{
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
