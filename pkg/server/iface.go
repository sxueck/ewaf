package server

import (
	_ "github.com/sxueck/ewaf/pkg/elog"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Name      string
	IPVersion string
	Listen    string
	Port      uint8
	IP        string
}

type StartServ interface {
	Start(Server) error // 承接请求
	Stop()              // 停止访问，优雅退出
	Serve() error       // 业务模块
}

type FrontendServer struct {
	StartServ
}

func (fs *FrontendServer) Start(s Server) error {
	e := echo.New()
	e.HideBanner = true
	logrus.Infof("Server listenner at IP: %s, Port %d, is starting", s.IP, s.Port)
	return nil
}

func (fs *FrontendServer) Stop() {
}

func (fs *FrontendServer) Serve() error {
	return nil
}
