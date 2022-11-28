package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
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

func Codec() encoding.Codec {
	return CodecWithParent(&protoCodec{})
}

type frame struct {
	payload []byte
}

func CodecWithParent(fallback encoding.Codec) encoding.Codec {
	return &rawCodec{fallback}
}

type rawCodec struct {
	parentCodec encoding.Codec
}

func (c *rawCodec) Name() string {
	return "proto"
}

func (c *rawCodec) Marshal(v interface{}) ([]byte, error) {
	out, ok := v.(*frame)
	if !ok {
		return c.parentCodec.Marshal(v)
	}
	return out.payload, nil

}

func (c *rawCodec) Unmarshal(data []byte, v interface{}) error {
	dst, ok := v.(*frame)
	if !ok {
		return c.parentCodec.Unmarshal(data, v)
	}
	dst.payload = data
	return nil
}

func (c *rawCodec) String() string {
	return fmt.Sprintf("proxy>%s", c.parentCodec.Name())
}

// protoCodec is a Codec implementation with protobuf
// it is the default rawCodec for gRPC.
type protoCodec struct{}

func (protoCodec) Name() string {
	return "proto"
}

func (protoCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (protoCodec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
