package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"path"
	"raspberrypi/protobuf"
	"raspberrypi/utils/dir"
	"raspberrypi/utils/execute"
	"raspberrypi/utils/file"
	"raspberrypi/utils/network"
)

type server struct {
	protobuf.UnimplementedInteractServer
}

func (s *server) SetAddr(ctx context.Context, setAddrRequest *protobuf.SetAddrRequest) (*protobuf.NormalResponse, error) {
	fmt.Printf("%s:%s\n", setAddrRequest.InterfaceName, setAddrRequest.InterfaceAddr)
	// ---------------- 核心逻辑 ----------------
	// 1. 设置非节能模式
	err := network.SetNoManagement(setAddrRequest.InterfaceName)
	if err != nil {
		return &protobuf.NormalResponse{
			Reply: "failed",
		}, fmt.Errorf("set no efficient failed: %v", err)
	}
	// 2. 进行地址的设置
	err = network.SetAddr(setAddrRequest.InterfaceName, setAddrRequest.InterfaceAddr, setAddrRequest.AddrType)
	if err != nil {
		fmt.Printf("set addr failed: %v\n", err)
		return &protobuf.NormalResponse{
			Reply: "failed",
		}, fmt.Errorf("set addr failed: %v", err)
	}
	// ---------------- 核心逻辑 ----------------
	return &protobuf.NormalResponse{
		Reply: "success",
	}, nil
}

func (s *server) AddRoute(ctx context.Context, addRouteRequest *protobuf.AddRouteRequest) (*protobuf.NormalResponse, error) {
	fmt.Printf("%s:%s\n", addRouteRequest.DestinationNetworkSegment, addRouteRequest.Gateway[:len(addRouteRequest.Gateway)-3])
	// ---------------- 核心逻辑 ----------------
	err := network.AddRoute(addRouteRequest.DestinationNetworkSegment, addRouteRequest.Gateway)
	if err != nil {
		fmt.Printf("add route failed: %v\n", err)
		return &protobuf.NormalResponse{
			Reply: "route already exists",
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
		}, fmt.Errorf("create directory failed: %v", err)
	}
	err = file.WriteStringIntoFile(transmitFileRequest.DestinationPath, transmitFileRequest.Content)
	if err != nil {
		fmt.Printf("write file failed: %v\n", err)
		return &protobuf.NormalResponse{
			Reply: "failed",
		}, fmt.Errorf("write file failed: %v", err)
	}
	fmt.Printf("write file %s success\n", transmitFileRequest.DestinationPath)
	// ---------------- 核心逻辑 ----------------
	return &protobuf.NormalResponse{
		Reply: "success",
	}, nil
}

func (s *server) SetEnv(ctx context.Context, setEnvRequest *protobuf.SetEnvRequest) (*protobuf.NormalResponse, error) {
	fmt.Println("set envs")
	finalString := ""
	// ---------------- 核心逻辑 ----------------
	for index, envField := range setEnvRequest.EnvFields {
		envValue := setEnvRequest.EnvValues[index]
		finalString += fmt.Sprintf("%s=%s\n", envField, envValue)
	}
	err := file.WriteStringIntoFile("/home/zeusnet/Projects/lir_node/lir_node/envs.txt", finalString)
	if err != nil {
		fmt.Printf("write envs file failed: %v\n", err)
		return &protobuf.NormalResponse{
			Reply: "failed",
		}, fmt.Errorf("write envs file failed: %v", err)
	}
	// ---------------- 核心逻辑 ----------------
	return &protobuf.NormalResponse{
		Reply: "success",
	}, nil
}

func (s *server) SetSysctls(ctx context.Context, setSysctlsRequest *protobuf.SetSysctlsRequest) (*protobuf.NormalResponse, error) {
	fmt.Println("set sysctls")
	// ---------------- 核心逻辑 ----------------
	for index, sysctlField := range setSysctlsRequest.SysctlFields {
		err := execute.Command("sysctl", []string{"-w", fmt.Sprintf("%s=%d", sysctlField, setSysctlsRequest.SysctlValues[index])})
		if err != nil {
			fmt.Printf("set sysctl %s failed: %v\n", sysctlField, err)
			return &protobuf.NormalResponse{
				Reply: "failed",
			}, fmt.Errorf("set sysctl %s failed: %v", sysctlField, err)
		}
	}
	// ---------------- 核心逻辑 ----------------
	return &protobuf.NormalResponse{
		Reply: "success",
	}, nil
}

func (s *server) LoadKernelInfo(ctx context.Context, loadKernelInfoRequest *protobuf.LoadKernelInfoRequest) (*protobuf.NormalResponse, error) {
	fmt.Println("load kernel info")
	// ---------------- 核心逻辑 ----------------
	// 调用 python 即可
	err := dir.WithContextManager("/home/zeusnet/Projects/lir_node/lir_node", func() error {
		err := execute.Command("/home/zeusnet/miniconda3/envs/lir/bin/python", []string{"start.py", "raspberrypi"})
		if err != nil {
			fmt.Printf("execute python failed: %v\n", err)
			return fmt.Errorf("execute python failed: %v", err)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("load kernel info failed: %v\n", err)
		return &protobuf.NormalResponse{
			Reply: "failed",
		}, fmt.Errorf("load kernel info failed: %v", err)
	}
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
