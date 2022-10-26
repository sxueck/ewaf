package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type StreamDirector func(context.Context, string) (context.Context, *grpc.ClientConn)

func NewProxy(dst *grpc.ClientConn, opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(
		append(opts,
			grpc.UnknownServiceHandler(TransparentHandler(DefaultDirector(dst))))...)
}

// DefaultDirector 流转发处理器
func DefaultDirector(cc *grpc.ClientConn) StreamDirector {
	return func(ctx context.Context, s string) (context.Context, *grpc.ClientConn) {
		md, _ := metadata.FromIncomingContext(ctx)
		ctx = metadata.NewOutgoingContext(ctx, md.Copy())
		return ctx, cc
	}
}
