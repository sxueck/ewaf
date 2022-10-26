package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// codec => stream director => stream handler => balancer

type handler struct {
	director StreamDirector
}

type ServerOptions struct {
	ctx context.Context
}

func (gso *ServerOptions) WithContext(ctx context.Context) {
	gso.ctx = ctx
}

func (gso *ServerOptions) Start() error {
	pBc := bufconn.Listen(10)

	// create a client connection to this backend
	//cc, err := backend

	return nil
}

func (gso *ServerOptions) Stop() {

}

func (gso *ServerOptions) Serve() error {
	return nil
}

func TransparentHandler(director StreamDirector) grpc.StreamHandler {
	streamer := &handler{
		director: director,
	}

	return streamer.handler
}

func (h *handler) handler(srv interface{}, serverStream grpc.ServerStream) error {
	return nil
}
