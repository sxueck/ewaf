package grpc

import (
	"github.com/sirupsen/logrus"
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
			logrus.Println(m)
			serviceName := m[:strings.LastIndexByte(m, '.')-1]
			methodName := m[len(serviceName):]
			RegisterService(gs, director, serviceName, methodName)
		}
	}
}

func backendDialer(opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	logrus.Printf("have connection")
	return nil, nil
}
