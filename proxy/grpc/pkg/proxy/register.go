package grpc

import (
	"google.golang.org/grpc"
)

func RegisterService(server *grpc.Server, director StreamDirector, serviceName string, methods ...string) {
	streamer := &handler{director}
	fakeDesc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*interface{})(nil),
	}

	for _, m := range methods {
		streamDesc := grpc.StreamDesc{
			StreamName:    m,
			Handler:       streamer.handler,
			ServerStreams: true,
			ClientStreams: true,
		}

		fakeDesc.Streams = append(fakeDesc.Streams, streamDesc)
	}

	server.RegisterService(fakeDesc, streamer)
}
