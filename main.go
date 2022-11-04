package main

import (
	"github.com/sirupsen/logrus"
	"github.com/sxueck/ewaf/proxy"
	"github.com/sxueck/ewaf/proxy/grpc"
	"golang.org/x/net/context"
)

func main() {
	//cfg := config.Cfg

	// global channel context
	ctx, _ := context.WithCancel(context.Background())

	// start internal proxy interfaces
	for _, f := range []proxy.StartServ{
		&grpc.ServerOptions{} /*, &http.ServerOptions{}, &tcp.ServerOptions{}*/} {
		go func(f proxy.StartServ) {
			f.WithContext(ctx)
			err := f.Start()
			if err != nil {
				logrus.Printf("fatal start internal server : %s", err)
				return
			}

			logrus.Println("start internal grpc serve module")

			err = f.Serve()
			if err != nil {
				logrus.Printf("fatal start internal server : %s", err)
				return
			}
		}(f)
	}

	<-ctx.Done()
}
