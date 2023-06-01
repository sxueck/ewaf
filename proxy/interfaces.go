package proxy

import (
	"context"
	"github.com/sxueck/ewaf/pkg"
	"log"
	"strings"
)

type StartServ interface {
	WithContext(context.Context, *pkg.GlobalConfig)

	Start() error
	Stop()
	Serve() error
}

func CheckTheSurvivalOfUpstreamServices(frs []pkg.Server, frsMark string) []pkg.Frontend {
	log.Printf("%+v\n", frs)
	tcpFrs := make([]pkg.Frontend, 0)
	for _, v := range frs {
		if strings.Contains(v.Frontend.Type, frsMark) {
			tcpFrs = append(tcpFrs, v.Frontend)
		}
	}
	return tcpFrs
}
