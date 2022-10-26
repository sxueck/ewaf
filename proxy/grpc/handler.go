package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"io"
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
	// a bit of grpc internals never hurt anyone
	methods, ok := grpc.MethodFromServerStream(serverStream)
	if !ok {
		return status.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
	}

	// require that director return context inherits from the serverStream context
	outCtx, backendConn := h.director(serverStream.Context(), methods)

	clientCtx, cancel := context.WithCancel(outCtx)
	defer cancel()

	globalStreamDesc := grpc.StreamDesc{
		ClientStreams: true,
		ServerStreams: true,
	}
	clientStream, err := grpc.NewClientStream(clientCtx, &globalStreamDesc, backendConn, methods)
	if err != nil {
		return err
	}

	// handler can switch net packet data
	s2cErrChan := h.forwardServerToClient(serverStream, clientStream)
	c2sErrChan := h.forwardClientToServer(clientStream, serverStream)

	// don't know which side is going to stop sending first,
	// so we need a select between the two
	for i := 0; i < 2; i++ {
		select {
		case s2cErr := <-s2cErrChan:
			if s2cErr == io.EOF {
				// the clientStream > serverStream may continue pumping though
				_ = clientStream.CloseSend()
			} else {
				// however, we may have gotten a reception error (stream disconnected, a read error etc.)
				// in which case we need to cancel the clientStream to the backend, let all of its goroutines
				// be freed up by the cancelfunc and exit with an error to the stack
				cancel()
				return status.Errorf(codes.Internal, "failed proxying s2c : %v", s2cErr)
			}
		case c2sErr := <-c2sErrChan:
			// this happens when the clientStream has nothing else to io.EOF
			// return a grpc error. In those two cases we may have received trailers as part of the call
			serverStream.SetTrailer(clientStream.Trailer())
			if c2sErr != io.EOF {
				return c2sErr
			}
			return nil
		}

	}

	return status.Errorf(codes.Internal, "grpc proxying should never reach this stage")
}

func (h *handler) forwardClientToServer(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	return nil
}

func (h *handler) forwardServerToClient(dst grpc.ServerStream, src grpc.ClientStream) chan error {
	return nil
}
