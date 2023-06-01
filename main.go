package main

import (
	"github.com/sirupsen/logrus"
	"github.com/sxueck/ewaf/config"
	"github.com/sxueck/ewaf/pkg"
	"github.com/sxueck/ewaf/proxy"
	"github.com/sxueck/ewaf/proxy/tcp"
	"golang.org/x/net/context"
)

func main() {
	cfg := config.InitParse(&pkg.GlobalConfig{})

	// global channel context
	ctx, _ := context.WithCancel(context.Background())

	// start internal proxy interfaces
	for _, f := range []proxy.StartServ{
		&tcp.ServerOptions{FrMark: "tcp"}} {
		go func(f proxy.StartServ) {
			f.WithContext(ctx, cfg.(*pkg.GlobalConfig))
			err := f.Start()
			if err != nil {
				logrus.Printf("fatal start internal server : %s", err)
				return
			}

			err = f.Serve()
			if err != nil {
				logrus.Printf("fatal start internal server : %s", err)
				return
			}
		}(f)
	}

	<-ctx.Done()
}
