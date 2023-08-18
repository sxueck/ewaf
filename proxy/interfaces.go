package proxy

import (
	"context"
	"github.com/sxueck/ewaf/pkg"
	"strings"
)

type StartServ interface {
	WithContext(context.Context, *pkg.GlobalConfig)

	Start() any
	Stop()
	Serve(any) error

	ExtraFrMark() string
}

func CheckTheSurvivalOfUpstreamServices(frs []pkg.Server, frsMark string) []pkg.Frontend {
	tcpFrs := make([]pkg.Frontend, 0)
	for _, v := range frs {
		if strings.Contains(v.Frontend.Type, frsMark) {
			tcpFrs = append(tcpFrs, v.Frontend)
		}
	}
	return tcpFrs
}
