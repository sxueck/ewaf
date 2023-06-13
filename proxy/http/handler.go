package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/sxueck/ewaf/pkg"
	"github.com/sxueck/ewaf/proxy"
	"net"
	"net/url"
	"os"
)

var tunName = "waf-tun0"

type ServerOptions struct {
	ctx    context.Context
	cfg    *pkg.GlobalConfig
	FrMark string // frontend mark

	lfpDir string
}

func (gso *ServerOptions) WithContext(ctx context.Context, cfg *pkg.GlobalConfig) {
	gso.ctx = ctx
	gso.cfg = cfg
}

func (gso *ServerOptions) ExtraFrMark() string {
	return gso.FrMark
}

func (gso *ServerOptions) Start() any {
	gso.lfpDir, _ = os.MkdirTemp("", tunName)
	return nil
}

func (gso *ServerOptions) Stop() {
	if len(gso.lfpDir) == 0 {
		return
	}
	if err := os.RemoveAll(gso.lfpDir); err != nil {
		logrus.Warn(err)
	}
}

func (gso *ServerOptions) Serve(in any) error {
	proxyTargets := proxy.CheckTheSurvivalOfUpstreamServices(gso.cfg.Servers, gso.FrMark)
	for _, v := range proxyTargets {
		var lis net.Listener
		go func(lis net.Listener) {
			var err error
			lis, err = net.Listen("unix", gso.lfpDir)
			if err != nil {
				logrus.Println(err)
				return
			}
		}(lis)

		e := echo.New()
		e.HideBanner = true
		e.Use(middleware.Recover())
		var targets []*middleware.ProxyTarget
		if len(v.Location) > 1 {
			logrus.Warn("HTTP multi-backend proxy is not supported for now")
			return nil
		}

		for _, vl := range v.Location {
			u := vl.Backend.ByPass
			urlParsing, err := url.Parse(u)
			if err != nil {
				logrus.Println("url parsing error, ", err)
				continue
			}

			targets = append(targets, &middleware.ProxyTarget{
				URL: urlParsing,
			})
		}

		e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
		e.Listener = lis
		err := e.Start("")
		if err != nil {
			return err
		}
	}
	return nil
}
