package http

import (
	"context"
	"github.com/sxueck/ewaf/pkg"
	"github.com/sxueck/ewaf/pkg/utils/water"
	"time"
)

var tunName = "waf-tun0"

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
	tun := water.NewTUNConfigure(tunName)
	err := tun.AssignTunChannelAddress(tun, gso.cfg.Global.TunChannelCidr)
	if err != nil {
		return err
	}
	<-time.NewTicker(3 * time.Minute).C
	return nil
}

func (gso *ServerOptions) Stop() {
	tun := water.NewTUNConfigure(tunName)
	tun.DeleteTUNChannel(tun)
}

func (gso *ServerOptions) Serve(in any) error {
	return nil
}
