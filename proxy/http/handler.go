package http

import (
	"context"
)

type ServerOptions struct {
	ctx context.Context
}

func (gso *ServerOptions) WithContext(ctx context.Context) {
	gso.ctx = ctx
}

func (gso *ServerOptions) Start() error {
	return nil
}

func (gso *ServerOptions) Stop() {

}

func (gso *ServerOptions) Serve() error {
	return nil
}
