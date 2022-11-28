// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: register.proto

package api_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ProxyRegisterServiceClient is the client API for ProxyRegisterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProxyRegisterServiceClient interface {
	// register rpc
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
}

type proxyRegisterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProxyRegisterServiceClient(cc grpc.ClientConnInterface) ProxyRegisterServiceClient {
	return &proxyRegisterServiceClient{cc}
}

func (c *proxyRegisterServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/ProxyRegisterService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProxyRegisterServiceServer is the server API for ProxyRegisterService service.
// All implementations should embed UnimplementedProxyRegisterServiceServer
// for forward compatibility
type ProxyRegisterServiceServer interface {
	// register rpc
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
}

// UnimplementedProxyRegisterServiceServer should be embedded to have forward compatible implementations.
type UnimplementedProxyRegisterServiceServer struct {
}

func (UnimplementedProxyRegisterServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}

// UnsafeProxyRegisterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProxyRegisterServiceServer will
// result in compilation errors.
type UnsafeProxyRegisterServiceServer interface {
	mustEmbedUnimplementedProxyRegisterServiceServer()
}

func RegisterProxyRegisterServiceServer(s grpc.ServiceRegistrar, srv ProxyRegisterServiceServer) {
	s.RegisterService(&ProxyRegisterService_ServiceDesc, srv)
}

func _ProxyRegisterService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyRegisterServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ProxyRegisterService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyRegisterServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProxyRegisterService_ServiceDesc is the grpc.ServiceDesc for ProxyRegisterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProxyRegisterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ProxyRegisterService",
	HandlerType: (*ProxyRegisterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _ProxyRegisterService_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "register.proto",
}
