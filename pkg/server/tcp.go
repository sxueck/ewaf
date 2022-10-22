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
			InterfaceName: "waf-tun0",
			Network:       "172.16.3.0/25",
		},
	}

	ifce, err := water.New(tunConfig)
	if err != nil {
		logrus.Fatal("CreateTUNChannel failed: ", err)
	}

	return ifce
}

// NewTcpStatusServer 每个httpserver前置一个tcp转发服务，用以控制 FIN/RST 等状态
func NewTcpStatusServer(port int) error {
	return nil
}
