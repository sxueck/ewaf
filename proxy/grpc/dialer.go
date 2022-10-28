package grpc

import (
	ecfg "github.com/sxueck/ewaf/config"
	"google.golang.org/grpc"
	"strings"
)

func GeneratorMethodsTree(gs *grpc.Server, sites ecfg.Cfg) {
	var director StreamDirector
	for _, s := range sites.Servers {
		if s.Frontend.Type != "grpc" {
			continue
		}
		for _, bk := range s.Frontend.Location {
			m := bk.Backend.Method
			serviceName := m[:strings.LastIndex(m, ":")-1]
			methodName := m[len(serviceName):]
			RegisterService(gs, director, serviceName, methodName)
		}
	}
}
