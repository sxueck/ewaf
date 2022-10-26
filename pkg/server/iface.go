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

type FrontendServer struct {
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
