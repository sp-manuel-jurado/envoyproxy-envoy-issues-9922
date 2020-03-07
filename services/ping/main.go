package main

import (
	"context"
	"net"

	pb "github.com/socialpoint/envoyproxy-envoy-issues-9922/services/ping/pkg/sp_rpc"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

const (
	addr = ":10005"
)

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPingServiceServer(grpcServer, serviceServer{})
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err.Error())
	}
}

type serviceServer struct {
}

func (s serviceServer) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
