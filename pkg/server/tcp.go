package server

import (
	_ "github.com/sxueck/ewaf/pkg/elog"

	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
)

func CreateTUNChannel() *water.Interface {
	tunConfig := water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name:        "waf-tun0",
			MultiQueue:  true,
			Permissions: &water.DevicePermissions{},
		},
	}

	iface, err := water.New(tunConfig)
	if err != nil {
		logrus.Fatal("CreateTUNChannel failed: ", err)
	}

	return iface
}

func CloseTUNChannel(iface *water.Interface) {

}

// NewTcpStatusServer 每个httpserver前置一个tcp转发服务，用以控制 FIN/RST 等状态
func NewTcpStatusServer(port int) error {

	return nil
}
