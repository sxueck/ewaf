package http

import (
	"context"
	"github.com/sxueck/ewaf/pkg"
	"github.com/sxueck/ewaf/pkg/server"
)

type ServerOptions struct {
	ctx    context.Context
	cfg    *pkg.GlobalConfig
	FrMark string // frontend mark
}

func (gso *ServerOptions) WithContext(ctx context.Context, cfg *pkg.GlobalConfig) {
	gso.ctx = ctx
	gso.cfg = cfg
}

func (gso *ServerOptions) Start() any {
	server.CreateTUNChannel()
	return nil
}

func (gso *ServerOptions) Stop() {

}

func (gso *ServerOptions) Serve(in any) error {
	return nil
}
