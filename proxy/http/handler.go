package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/sxueck/ewaf/pkg"
	"github.com/sxueck/ewaf/pkg/utils/water"
	"github.com/sxueck/ewaf/proxy"
	"os"
	"reflect"
)

var tunName = "waf-tun0"

type ServerOptions struct {
	ctx    context.Context
	cfg    *pkg.GlobalConfig
	FrMark string // frontend mark

	lfp *os.File
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
	gso.lfp = tun.File()
	return nil
}

func (gso *ServerOptions) Stop() {
	tun := water.NewTUNConfigure(tunName)
	tun.DeleteTUNChannel(tun)
}

func (gso *ServerOptions) Serve(in any) error {
	if reflect.DeepEqual(gso.lfp, os.File{}) {
		logrus.Fatalln("listening channel does not exist")
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	proxyTargets := proxy.CheckTheSurvivalOfUpstreamServices(gso.cfg.Servers, gso.FrMark)
	return nil
}
