package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func backendDialer() (*grpc.ClientConn, error) {
	backendBc := bufconn.Listen(10)
	srv := grpc.NewServer()

}
